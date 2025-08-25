package logger

import (
	"log"
	"os"
)

var (
	InfoLogger   *log.Logger
	ErrorLogger  *log.Logger
	SystemLogger *log.Logger
	file         *os.File
)

func init() {
	filename := "app.log"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	SystemLogger = log.New(file, "SYSTEM: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Close() {
	if file != nil {
		file.Close()
	}
}

func Info(msg string) {
	InfoLogger.Output(2, msg)
}
func Error(msg string) {
	ErrorLogger.Output(2, msg)
}
func System(msg string) {
	SystemLogger.Output(2, msg)
}
