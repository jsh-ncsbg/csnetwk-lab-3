package logger

import (
	"fmt"
	"io"
	"time"
)

type Logger interface {
	Println(...any)
	Printf(string, ...any)
}

type logger struct {
	Out  io.Writer
	name string
}

func (l *logger) Println(v ...any) {
	fmt.Fprintln(l.Out, l.makePrefix(), fmt.Sprint(v...))
}

func (l *logger) Printf(format string, v ...any) {
	fullFormat := l.makePrefix() + " " + format

	fmt.Fprintf(l.Out, fullFormat, v...)
}

func (l *logger) makePrefix() string {
	return time.Now().Format(time.TimeOnly) + " [" + l.name + "]"
}

func NewLogger(name string, out io.Writer) *logger {
	return &logger{
		Out:  out,
		name: name,
	}
}
