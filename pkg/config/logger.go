package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func LoggerLoad() {
	// set logger location open a file
	logFile, err := os.OpenFile(Config.LogPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(logFile)

	log.SetReportCaller(true)

	// Only log the warning severity or above.
	log.SetLevel(log.TraceLevel)

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
}
