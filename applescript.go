package thyme 

const (
active_windows_script = `on logWindow(proc, each_window)
  try
    log "WINDOW " & (id of each_window) & ":" & (name of each_window) as string
  on error e
    log "WINDOW " & (id of proc) & ":" & (name of each_window) as string
  end try
end logWindow

set nonElectronProcesses to {}
tell application "System Events"
  repeat with proc in (first application process whose frontmost is true)
    set procName to (name of proc)
    if procName is "Electron" then
      set app_windows to item 1 of (every window of proc)
      repeat with each_window in app_windows
        my logWindow(proc, each_window)
      end repeat
    else
      set end of nonElectronProcesses to proc
    end if
  end repeat
end tell

repeat with proc in nonElectronProcesses
  set app_windows to item 1 of (every window of proc)
  repeat with each_window in app_windows
    my logWindow(proc, each_window)
  end repeat
end repeat
`
all_windows_script = `on logWindow(proc, each_window)
  try
    log "WINDOW " & (id of each_window) & ":" & (name of each_window) as string
  on error e
    log "WINDOW " & (id of proc) & ":" & (name of each_window) as string
  end try
end logWindow

set nonElectronProcesses to {}
tell application "System Events"
  repeat with proc in (every application process where background only is false)
    set procName to (name of proc)
    if procName is "Electron" then
      log "PROCESS " & (id of proc) & ":" & procName
      set app_windows to (every window of proc)
      repeat with each_window in app_windows
        my logWindow(proc, each_window)
      end repeat
    else
      set end of nonElectronProcesses to proc
    end if
  end repeat
end tell

repeat with proc in nonElectronProcesses
  log "PROCESS " & (id of proc) & ":" & (name of proc)
  set app_windows to (every window of proc)
  repeat with each_window in app_windows
    my logWindow(proc, each_window)
  end repeat
end repeat
`
visible_windows_script = `on logWindow(proc, each_window)
  try
    log "WINDOW " & (id of each_window) & ":" & (name of each_window) as string
  on error e
    log "WINDOW " & (id of proc) & ":" & (name of each_window) as string
  end try
end logWindow

set nonElectronProcesses to {}
tell application "System Events"
  repeat with proc in (every process whose visible is true)
    set procName to (name of proc)
    if procName is "Electron" then
      log "PROCESS " & (id of proc) & ":" & procName
      set app_windows to (every window of proc)
      repeat with each_window in app_windows
        my logWindow(proc, each_window)
      end repeat
    else
      set end of nonElectronProcesses to proc
    end if
  end repeat
end tell

repeat with proc in nonElectronProcesses
  log "PROCESS " & (id of proc) & ":" & (name of proc)
  set app_windows to (every window of proc)
  repeat with each_window in app_windows
    my logWindow(proc, each_window)
  end repeat
end repeat
`
)
