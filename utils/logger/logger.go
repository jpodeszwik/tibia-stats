package logger

import (
	"io"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
	Debug *log.Logger

	infoOutput  = os.Stdout
	errorOutput = os.Stderr
	debugOutput = io.Discard
)

const infoPrefix = "INFO: "
const errorPrefix = "ERROR: "

func init() {
	logfile, exists := os.LookupEnv("LOGFILE")
	if exists {
		file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file %v", err)
		}

		infoOutput = file
		errorOutput = file
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "DEBUG" {
		debugOutput = infoOutput
	}

	Info = log.New(infoOutput, infoPrefix, log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorOutput, errorPrefix, log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(debugOutput, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}
