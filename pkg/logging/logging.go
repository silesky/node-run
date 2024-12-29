package logging

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
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// SetDebug sets the debug flag
func SetDebug(d bool) {
	debug = d
}

// Info logs an info message if debugging is enabled
func Info(v ...interface{}) {
	if debug {
		InfoLogger.Println(v...)
	}
}

// Infof logs a formatted info message if debugging is enabled
func Infof(format string, v ...interface{}) {
	if debug {
		InfoLogger.Printf(format, v...)
	}
}

// Error logs an error message
func Error(v ...interface{}) {
	ErrorLogger.Println(v...)
}

// Errorf logs a formatted error message
func Errorf(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}
