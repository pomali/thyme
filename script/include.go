package main

import (
  "io"
  "io/ioutil"
  "os"
  "strings"
)

// Reads all .applescript files in the current folder
// and encodes them as strings literals in applescript.go
func main() {
  fs, _ := ioutil.ReadDir("./script")
  out, _ := os.Create("applescript.go")
  out.Write([]byte("package thyme \n\nconst (\n"))
  for _, f := range fs {
    if strings.HasSuffix(f.Name(), ".applescript") {
      out.Write([]byte(strings.TrimSuffix(f.Name(), ".applescript") + " = `"))
      f, _ := os.Open("./script/" + f.Name())
      io.Copy(out, f)
      out.Write([]byte("`\n"))
    }
  }
  out.Write([]byte(")\n"))
}
