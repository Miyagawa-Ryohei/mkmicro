package infra

import (
	"fmt"
	"github.com/Miyagawa-Ryohei/mkmicro/types"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type logger struct {
	buf chan string
	wg  *sync.WaitGroup
}

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func (l *logger) getCaller() (filename string, funcName string, line int) {
	filename = ""
	pt, file, line, ok := runtime.Caller(3)
	if ok {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		} else {
			filename = strings.Replace(file, cwd+"/", "", -1)
		}
	}
	funcName = runtime.FuncForPC(pt).Name()
	return
}

func (l *logger) getColorFormat(logLevel LogLevel) string {
	switch logLevel {
	case DEBUG:
		return "\x1b[32m"
	case INFO:
		return "\x1b[34m"
	case WARN:
		return "\x1b[33m"
	case ERROR:
		return "\x1b[31m"
	default:
		return "\x1b[30m"
	}
}

func (l *logger) Print(logLevel LogLevel, msg string) {
	filename, funcName, line := l.getCaller()
	color := l.getColorFormat(logLevel)
	tag := os.Getenv("LOG_TAG")
	if tag == "" {
		tag = "no tag"
	}
	l.buf <- fmt.Sprintf("%s[%s][%s][%s:%d][%s][%s] %s\x1b[0m", color, tag, time.Now().Format(time.RFC3339), filename, line, funcName, logLevel.String(), msg)
}

func (l *logger) Info(msg string) {
	l.Print(INFO, msg)
}

func (l *logger) Infof(format string, binder ...interface{}) {
	l.Print(INFO, fmt.Sprintf(format, binder...))
}

func (l *logger) Debug(msg string) {
	l.Print(DEBUG, msg)
}

func (l *logger) Debugf(format string, binder ...interface{}) {
	l.Print(DEBUG, fmt.Sprintf(format, binder...))
}

func (l *logger) Warn(msg string) {
	l.Print(WARN, msg)
}

func (l *logger) Warnf(format string, binder ...interface{}) {
	l.Print(WARN, fmt.Sprintf(format, binder...))
}

func (l *logger) Error(msg string) {
	l.Print(ERROR, msg)
}

func (l *logger) Errorf(format string, binder ...interface{}) {
	l.Print(ERROR, fmt.Sprintf(format, binder...))
}

func (l *logger) Flush() {
	l.wg.Add(1)
	l.buf <- ""
	l.wg.Wait()
}

func (l *logger) printLog() {
	for {
		msg, more := <-l.buf
		if more {
			if msg == "" {
				l.wg.Done()
			} else {
				fmt.Println(msg)
			}
		} else {
			return
		}
	}
}

func NewLogger() types.Logger {

	buf := make(chan string)
	wg := &sync.WaitGroup{}

	l := &logger{
		wg:  wg,
		buf: buf,
	}
	go l.printLog()
	return l
}

var DefaultLogger = NewLogger()
