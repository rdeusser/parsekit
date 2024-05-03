package logging

import (
	"fmt"
	"io"
	"time"
)

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Error
)

type Logger struct {
	w io.Writer
	l LogLevel
}

func New(w io.Writer, l LogLevel) Logger {
	return Logger{
		w: w,
		l: l,
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.l = level
}

func (l Logger) Debug(msg string, args ...any) {
	if l.l <= Debug {
		fmt.Fprintln(l.w, format(Debug, msg, args...))
	}
}

func (l Logger) Info(msg string, args ...any) {
	if l.l <= Info {
		fmt.Fprintln(l.w, format(Info, msg, args...))
	}
}

func (l Logger) Error(msg string, args ...any) {
	if l.l <= Error {
		fmt.Fprintln(l.w, format(Error, msg, args...))
	}
}

func format(l LogLevel, msg string, args ...any) string {
	now := time.Now()
	level := ""
	switch l {
	case Debug:
		level = "DEBUG"
	case Info:
		level = "INFO "
	case Error:
		level = "ERROR"
	}
	return fmt.Sprintf("%s %s %s %s", now.Format(time.DateOnly), now.Format(time.Kitchen), level, fmt.Sprintf(msg, args...))
}
