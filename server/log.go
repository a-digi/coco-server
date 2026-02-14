package server

import (
	"fmt"
	"time"
)

// LogFileName returns a log file name based on the current datetime (YYYYMMDD_HHMMSS.log)
func LogFileName(prefix string) string {
	return fmt.Sprintf("%s_%s.log", prefix, time.Now().Format("20060102_150405"))
}
