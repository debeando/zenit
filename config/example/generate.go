//go:build ignore
// +build ignore

package main

import (
	"io"
	"os"
)

func main() {
	in, err := os.Open("../../zenit.yaml")
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create("example.go")
	if err != nil {
		return
	}
	defer out.Close()

	out.Write([]byte("package example\n\nconst File = `"))
	io.Copy(out, in)
	out.Write([]byte("`\n"))
}
