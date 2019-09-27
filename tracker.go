package thyme

import "log"

// Tracker tracks application usage. An implementation that satisfies
// this interface is required for each OS windowing system Thyme
// supports.
type Tracker interface {
	// Snap returns a Snapshot reflecting the currently in-use windows
	// at the current time.
	Snap() (*Snapshot, error)

	// Deps returns a string listing the dependencies that still need
	// to be installed with instructions for how to install them.
	Deps() string
}

// trackers is the list of Tracker constructors that are available on this system. Tracker implementations should call
// the RegisterTracker function to make themselves available.
var trackers = make(map[string]func() Tracker)

// RegisterTracker makes a Tracker constructor available to clients of this package.
func RegisterTracker(name string, t func() Tracker) {
	if _, exists := trackers[name]; exists {
		log.Fatalf("a tracker already exists with the name %s", name)
	}
	trackers[name] = t
}

// NewTracker returns a new Tracker instance whose type is `name`.
func NewTracker(name string) Tracker {
	if _, exists := trackers[name]; !exists {
		log.Fatalf("no Tracker constructor has been registered with name %s", name)
	}
	return trackers[name]()
}
