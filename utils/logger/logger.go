package logger

import (
	"io"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

const infoPrefix = "INFO: "
const errorPrefix = "ERROR: "

func init() {
	logfile, exists := os.LookupEnv("LOGFILE")
	var output io.Writer
	if exists {
		var err error
		output, err = os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file %v", err)
		}
		Info = log.New(output, infoPrefix, log.Ldate|log.Ltime|log.Lshortfile)
		Error = log.New(output, errorPrefix, log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		Info = log.New(os.Stdout, infoPrefix, log.Ldate|log.Ltime|log.Lshortfile)
		Error = log.New(os.Stderr, errorPrefix, log.Ldate|log.Ltime|log.Lshortfile)
	}
}
