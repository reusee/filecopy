package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/reusee/e/v2"
)

var (
	me     = e.Default.WithStack()
	ce, he = e.New(me)
	pt     = fmt.Printf
)

func main() {
	var from string
	flag.StringVar(&from, "from", "", "from file path")
	var to string
	flag.StringVar(&to, "to", "", "to file path")
	flag.Parse()

	if from == "" {
		panic("no from file")
	}
	if to == "" {
		panic("no to file")
	}

	fromFile, err := os.OpenFile(from, os.O_RDONLY, 0)
	ce(err, "open from file")
	defer fromFile.Close()
	toFile, err := os.OpenFile(to, os.O_WRONLY, 0)
	ce(err, "open to file")
	defer toFile.Close()

	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * time.Duration(i))
		fmt.Printf("ALL DATA IN %s WILL LOST. IS IT OK? (yes/no) ", to)
		os.Stdout.Sync()
		var confirm string
		fmt.Scanf("%s", &confirm)
		if confirm != "yes" {
			fmt.Printf("do nothing\n")
			return
		}
	}

	var bytesCopied int64
	go func() {
		for range time.NewTicker(time.Second).C {
			pt("%s copied\n", formatBytes(atomic.LoadInt64(&bytesCopied)))
		}
	}()

	buf := make([]byte, 256*1024*1024)
	for {
		n, err := fromFile.Read(buf)
		if n > 0 {
			_, e := toFile.Write(buf[:n])
			ce(e)
			atomic.AddInt64(&bytesCopied, int64(n))
		}
		if errors.Is(err, io.EOF) {
			break
		}
		ce(err)
	}

}

var units = []string{"B", "K", "M", "G", "T", "P", "E", "Z", "Y"}

func formatBytes(n int64) string {
	if n == 0 {
		return "0"
	}
	return strings.TrimSpace(_formatBytes(n, 0))
}

func _formatBytes(n int64, unitIndex int) string {
	if n == 0 {
		return ""
	}
	var str string
	next := n / 1024
	rem := n - next*1024
	if rem > 0 {
		str = fmt.Sprintf(" %d%s", rem, units[unitIndex])
	}
	return _formatBytes(next, unitIndex+1) + str
}
