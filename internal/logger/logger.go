package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zeroaddresss/golang-unisat-monitor/internal/utils"
)

type Logger struct {
	LogFile    *os.File
	InfoLogger *log.Logger
	ErrLogger  *log.Logger
}

var L *Logger

func NewLogger(dir string) *Logger {
	t := time.Now()
	timestamp := t.Format("02-01-2006_15:04")
	filePath := fmt.Sprintf("%s%s.txt", dir, timestamp)

	// create logs directory if it doesn't exist
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(file, "[INFO]  ", log.Ldate|log.Ltime)
	errLog := log.New(file, "[ERROR] ", log.Ldate|log.Ltime)

	return &Logger{
		LogFile:    file,
		InfoLogger: logger,
		ErrLogger:  errLog,
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.InfoLogger.Printf(format, v...)
	log.Printf(utils.Yellow("[INFO]  → ")+format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.ErrLogger.Printf(format, v...)
	log.Printf(utils.Red("[ERROR] → ")+format, v...)
}

func FormatFloat(f float64) string {
	if f == float64(int64(f)) {
		return fmt.Sprintf("%d", int64(f))
	}
	return fmt.Sprintf("%.3f", f)
}
