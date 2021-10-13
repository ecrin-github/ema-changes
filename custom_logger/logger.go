package custom_logger

import (
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	filename      string
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
}

var Clogger *Logger
var once sync.Once

func GetInstance(fname string) *Logger {
	once.Do(func() {
		Clogger = createLogger(fname)
	})
	return Clogger
}

func createLogger(fname string) *Logger {
	date_stamp := time.Now().Format("2006-02-01")
	file_full_name := fname + " " + date_stamp
	file, _ := os.OpenFile(file_full_name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	return &Logger{
		filename:      file_full_name,
		InfoLogger:    log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		WarningLogger: log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLogger:   log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
