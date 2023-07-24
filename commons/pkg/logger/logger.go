package logger

import (
	"log"
	"os"
)

var lInfo = log.New(os.Stderr, "[INFO]", log.Ldate|log.Ltime|log.Lshortfile)
var lWarn = log.New(os.Stderr, "[WARN]", log.Ldate|log.Ltime|log.Lshortfile)
var lError = log.New(os.Stderr, "[ERROR]", log.Ldate|log.Ltime|log.Lshortfile)

func Info(format string, v ...any) {
	lInfo.Printf(format, v...)
}

func Warn(format string, v ...any) {
	lWarn.Printf(format, v...)
}

func Error(format string, v ...any) {
	lError.Printf(format, v...)
}