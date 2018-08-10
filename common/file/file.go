package file

import (
  "os"
  "io/ioutil"
)

func Create(s string) bool {
  // detect if file exists
  var _, err = os.Stat(s)

  // create file if not exists
  if os.IsNotExist(err) {
    var file, err = os.Create(s)
    if err != nil { return false }
    defer file.Close()
  }

  return true
}

func Write(f string, s string) bool {
  // open file using READ & WRITE permission
  file, err := os.OpenFile(f, os.O_RDWR, 0644)
  defer file.Close()
  if err != nil { return false }
  
  // write some text line-by-line to file
  _, err = file.WriteString(s)
  if err != nil { return false }

  // save changes
  err = file.Sync()
  if err != nil { return false }

  return true
}

func Read(path string) (lines string) {
  if _, err := os.Stat(path); err == nil {
    contents, err := ioutil.ReadFile(path)
    if err != nil {
      return
    }

    lines = string(contents)
  }
  return
}

func Truncate(f string) bool {
  // open file using READ & WRITE permission
  file, err := os.OpenFile(f, os.O_RDWR, 0644)
  defer file.Close()
  if err != nil { return false }
  
  file.Truncate(0)
  file.Seek(0,0)
  file.Sync()

  return true
}
