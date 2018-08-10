package file

import (
  "os"
  "io/ioutil"
)

func Exist(f string) bool {
  if _, err := os.Stat(f); err != nil {
    return false
  }
  return true
}

func Create(f string) bool {
  if ! Exist(f) {
    var file, err = os.Create(f)
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

func Read(f string) (lines string) {
  if _, err := os.Stat(f); err == nil {
    contents, err := ioutil.ReadFile(f)
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

func Delete(f string) bool {
  if err := os.Remove(f); err != nil {
    return false
  }
  return true
}
