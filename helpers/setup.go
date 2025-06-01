package helpers

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// SetupLogger sets up the log output to a timestamped file inside the "./log" folder.
// 	- It creates the "./log" folder if it does not exist.
func SetupLogger() (*os.File, error) {
	logDir := "./log"

	// Check if log directory exists, create if it doesn't
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	logFilename := "logfile_" + time.Now().Format("20060102_150405") + ".log"
	logFilePath := filepath.Join(logDir, logFilename)

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	return file, nil
}
