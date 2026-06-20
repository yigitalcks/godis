package logging

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const timestampFormat = "2006-01-02 15:04:05.000"

type Logger struct {
	out io.Writer
	mu  sync.Mutex
}

var defaultLogger = &Logger{out: os.Stdout}

func (l *Logger) prefix() string {
	return time.Now().Format(timestampFormat)
}

func (l *Logger) Println(args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintln(l.out, append([]any{l.prefix()}, args...)...)
}

func (l *Logger) Printf(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.out, "%s "+format, append([]any{l.prefix()}, args...)...)
}

func Println(args ...any) {
	defaultLogger.Println(args...)
}

func Printf(format string, args ...any) {
	defaultLogger.Printf(format, args...)
}
