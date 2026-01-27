/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package game

import (
	"io"
	"time"
)

type writer struct {
	io.Writer
	timeFormat string
}

func (w writer) Write(b []byte) (n int, err error) {
	return w.Writer.Write(append([]byte(time.Now().Format(w.timeFormat)), b...))
}

// func yellow(s string) string {
// 	return "\033[33m" + s + "\033[0m"
// }

func red(s string) string {
	return "\033[31m" + s + "\033[0m"
}
