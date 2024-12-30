package logger

import (
	"log"
	"os"
)

var (
	DebugLogger *log.Logger
	debug       = false
)

func init() {
	DebugLogger = log.New(os.Stdout, "Debug: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// SetDebug sets the debug flag
func SetDebug(d bool) {
	debug = d
}

// Debug logs an info message if debugging is enabled
func Debug(v ...interface{}) {
	if debug {
		DebugLogger.Println(v...)
	}
}

// Debugf logs a formatted info message if debugging is enabled
func Debugf(format string, v ...interface{}) {
	if debug {
		DebugLogger.Printf(format, v...)
	}
}
