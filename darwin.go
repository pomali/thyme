package thyme

import (
	"bytes"
	"fmt"
	"hash/fnv"
  //"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//go:generate go run script/include.go

var (
	// allWindowsScript fetches the windows of all scriptable applications.  It
	// iterates through each application process known to System Events and attempts
	// to script the application with the same name as the application process. If
	// such an application exists and is scriptable, it prints the name of every
	// window in the application. Otherwise, it just prints the name of every
	// visible window in the application. If no visible windows exist, it will just
	// print the application name.  (System Events processes only have windows in
	// the current desktop/workspace.)
	allWindowsScript    string
	activeWindowsScript string
	// visibleWindowsScript generates a mapping from process to windows in the
	// current desktop (note: this is slightly different than the behavior of
	// the previous two scripts, where an empty windows list for a process
	// should NOT imply that there is one window named after the process.
	// Furthermore, the window IDs are not valid in this script (only the window
	// name is valid).
	visibleWindowsScript string
)

func init() {
	RegisterTracker("darwin", NewDarwinTracker)
	allWindowsScript = all_windows_script
	activeWindowsScript = active_windows_script
	visibleWindowsScript = visible_windows_script
}

// DarwinTracker tracks application usage using the "System Events" API in AppleScript. Due to the liminations of this
// API, the DarwinTracker will not be able to detect individual windows of applications that are not scriptable (in the
// AppleScript sense). For these applications, a single window is emitted with the name set to the application process
// name and the ID set to the process ID.
type DarwinTracker struct{}

var _ Tracker = (*DarwinTracker)(nil)

func NewDarwinTracker() Tracker {
	return &DarwinTracker{}
}

func (t *DarwinTracker) Deps() string {
	return `
  v0.1
		You will need to enable privileges for "Terminal" in System Preferences > Security & Privacy > Privacy > Accessibility.
		See https://support.apple.com/en-us/HT202802 for details.

		Note: this command prints out this message regardless of whether this has been done or not.
		`
}

func (t *DarwinTracker) Snap() (*Snapshot, error) {
	var allWindows []*Window
	var allProcWins map[process][]*Window
	{
		procWins, err := runAS(allWindowsScript)
		if err != nil {
			return nil, err
		}
		for proc, wins := range procWins {
			if len(wins) == 0 {
				allWindows = append(allWindows, &Window{ID: proc.id, Name: proc.name})
			} else {
				allWindows = append(allWindows, wins...)
			}
		}
		allProcWins = procWins
	}

	var active int64
	{
		procWins, err := runAS(activeWindowsScript)
		if err != nil {
			return nil, err
		}
		if len(procWins) > 1 {
			return nil, fmt.Errorf("found more than one active process: %+v", procWins)
		}
		for proc, wins := range procWins {
			if len(wins) == 0 {
				active = proc.id
			} else if len(wins) == 1 {
				active = wins[0].ID
			} else {
				return nil, fmt.Errorf("found more than one active window: %+v", wins)
			}
		}
	}

	var visible []int64
	{
		procWins, err := runAS(visibleWindowsScript)
		if err != nil {
			return nil, err
		}
		for proc, wins := range procWins {
			allWins := allProcWins[proc]
			for _, visWin := range wins {
				if len(allWins) == 0 {
					visible = append(visible, proc.id)
				} else {
					found := false
					for _, win := range allWins {
						if win.Name == visWin.Name {
							visible = append(visible, win.ID)
							found = true
							break
						}
					}
					if !found {
						log.Printf("warning: window ID not found for visible window %q", visWin.Name)
					}
				}
			}
		}
	}

	return &Snapshot{
		Time:    time.Now(),
		Windows: allWindows,
		Active:  active,
		Visible: visible,
	}, nil
}

// process is the {name, id} of a process
type process struct {
	name string
	id   int64
}

// runAS runs script as AppleScript and parses the output into a map of
// processes to windows.
func runAS(script string) (map[process][]*Window, error) {
	cmd := exec.Command("osascript")
	cmd.Stdin = bytes.NewBuffer([]byte(script))
	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("AppleScript error: %s, output was:\n%s", err, string(b))
	}
	return parseASOutput(string(b))
}

// parseASOutput parses the output of the AppleScript snippets used to extract window information.
func parseASOutput(out string) (map[process][]*Window, error) {
	proc := process{}
	procWins := make(map[process][]*Window)
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "PROCESS ") {
			c := strings.Index(line, ":")
			procID, err := strconv.ParseInt(line[len("PROCESS "):c], 10, 0)
			if err != nil {
				return nil, err
			}
			proc = process{line[c+1:], procID}
			procWins[proc] = nil
		} else if strings.HasPrefix(line, "WINDOW ") {
			win, winID := parseWindowLine(line, proc.id)
			procWins[proc] = append(procWins[proc],
				&Window{ID: winID, Name: fmt.Sprintf("%s - %s", win, proc.name)},
			)
		}
	}
	return procWins, nil
}

// parseWindowLine parses window ID from a line of the AppleScript
// output. If the ID is missing ("missing value"), parseWindowLine
// will return the hash of the window title and process ID. Note: if 2
// windows controlled by the same process both have IDs missing and
// have the same title, they will hash to the same ID. This is
// unfortunate but seems to be the best behavior.
func parseWindowLine(line string, procId int64) (string, int64) {
	c := strings.Index(line, ":")
	win := line[c+1:]
	winID, err := strconv.ParseInt(line[len("WINDOW "):c], 10, 0)
	if err != nil {
		// sometimes "missing value" appears here, so generate a value
		// taking the process ID and the window index to generate a hash
		winID = hash(fmt.Sprintf("%s%v", win, procId))
	}
	return win, winID
}


// hash converts a string to an integer hash
func hash(s string) int64 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int64(h.Sum32())
}
