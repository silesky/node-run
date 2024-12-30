package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	debug       bool
)

func init() {
	InfoLogger = log.New(os.Stdout, "Debug: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// SetDebug sets the debug flag
func SetDebug(d bool) {
	debug = d
}

// Debug logs an info message if debugging is enabled
func Debug(v ...interface{}) {
	if debug {
		InfoLogger.Println(v...)
	}
}

// Debugf logs a formatted info message if debugging is enabled
func Debugf(format string, v ...interface{}) {
	if debug {
		InfoLogger.Printf(format, v...)
	}
}
