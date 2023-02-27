package common

import (
	"os"

	"github.com/hpcloud/tail"
)

func Tail(filename string, out chan<- string) {
	defer close(out)

	t, _ := tail.TailFile(filename, tail.Config{
		Follow: true,
		Logger: tail.DiscardingLogger,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: os.SEEK_END,
		},
	})
	for line := range t.Lines {
		out <- line.Text
	}
}
