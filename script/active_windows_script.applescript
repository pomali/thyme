tell application "System Events"
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

