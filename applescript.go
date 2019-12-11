package thyme 

const (
active_windows_script = `tell application "System Events"
	set proc to (first application process whose frontmost is true)
	set procName to (name of proc)
	try
		tell application procName
			log "WINDOW " & (id of window 1) & ":" & (name of window 1)
		end tell
	on error e
		log "WINDOW " & (id of proc) & ":" & (name of first window of proc)
	end try
end tell

`
all_windows_script = `tell application "System Events"
	repeat with proc in (every application process where background only is false)
		tell process proc
			log "PROCESS " & (id of proc) & ":" & (name of proc)
			repeat with w in (windows of proc)
				log "WINDOW -:" & (name of w) as string
			end repeat
		end tell
	end repeat
end tell`
visible_windows_script = `tell application "System Events"
	set listOfProcesses to (every process whose visible is true)
	
	repeat with proc in listOfProcesses
		set procName to (name of proc)
		set procID to (id of proc)
		log "PROCESS " & procID & ":" & procName
		set app_windows to (every window of proc)
		repeat with each_window in app_windows
			log "WINDOW -:" & (name of each_window) as string
		end repeat
	end repeat
	
end tell`
)
