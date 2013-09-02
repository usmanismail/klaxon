package log

import (
   	"log"
)

func Info(logMessage string) {
	log.Printf("[INFO] "+logMessage) 
}

func Debug(logMessage string) {
	log.Printf("[DEBUG] "+logMessage) 
}

func Warn(logMessage string) {
	log.Printf("[WARN] "+logMessage) 
}

func Error(logMessage string) {
	log.Printf("[ERROR] "+logMessage) 
}