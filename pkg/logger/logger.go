package logger

import (
	"log"
	"os"
)

var (
	DebugLogger = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	debug       = false
)

func SetDebug(d bool) {
	debug = d
}

func Debug(v ...interface{}) {
	if debug {
		DebugLogger.Println(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if debug {
		DebugLogger.Printf(format, v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	DebugLogger.Fatalf(format, v...)
}
