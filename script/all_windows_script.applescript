tell application "System Events"
  set listOfProcesses to (every application process where background only is false)
end tell

repeat with proc in listOfProcesses
  set procName to (name of proc)
  set procID to (id of proc)
  log "PROCESS " & procID & ":" & procName
  -- Attempt to list windows if the process is scriptable
  try
    tell application procName
      repeat with i from 1 to (count windows)
        log "WINDOW " & (id of window i) & ":" & (name of window i) as string
      end repeat
    end tell
    on error err
      log err
  end try
end repeat
