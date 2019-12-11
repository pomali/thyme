tell application "System Events"
	repeat with proc in (every application process where background only is false)
		tell process proc
			log "PROCESS " & (id of proc) & ":" & (name of proc)
			repeat with w in (windows of proc)
				log "WINDOW -:" & (name of w) as string
			end repeat
		end tell
	end repeat
end tell