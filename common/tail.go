package common

import (
  "io"
  "os"
  "time"
)

func Tail(filename string, out chan<- string) error {
  if _, err := os.Stat(filename); os.IsNotExist(err) {
    return err
  }

  file, err := os.Open(filename)
  if err != nil {
    return err
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
        break
      }
    }

    offset += int64(readBytes)
    if readBytes != 0 {
      out <- string(buffer[:readBytes])
    }
    time.Sleep(250 * time.Millisecond)
  }
  return nil
}
