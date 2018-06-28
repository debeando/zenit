package main

import (
  "fmt"
  "io"
  "os"
  "time"
)

func main() {
  filename := os.Args[1]
  if _, err := os.Stat(filename); os.IsNotExist(err) {
    fmt.Printf("No exist file...")
  }

  file, err := os.Open(filename)
  if err != nil {
    fmt.Println("Failed to open file: ", filename)
    return
  }
  defer file.Close()

  offset, _ := file.Seek(0, io.SeekEnd)
  offset_tmp, _ := file.Seek(0, io.SeekEnd)
  buffer := make([]byte, 1024, 1024)

  for {
    offset_tmp, _ = file.Seek(0, io.SeekEnd)
    if offset_tmp < offset {
      offset, _ = file.Seek(0, io.SeekStart)
    }

    readBytes, err := file.ReadAt(buffer, offset)
    if err != nil {
      if err != io.EOF {
        fmt.Println("Error reading lines:", err)
        break
      }
    }

    offset += int64(readBytes)
    if readBytes != 0 {
      s := string(buffer[:readBytes])
      fmt.Printf("%s",s)
    }
    time.Sleep(250 * time.Millisecond)
  }
}
