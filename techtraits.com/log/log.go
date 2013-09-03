package log

import (
   	"log"
)

func Info(logMessage string, v ...interface{}) {
	if v != nil {
		log.Printf("[INFO]  "+ logMessage, v)
	} else {
		log.Printf("[INFO]  "+ logMessage)
	} 
}

func Debug(logMessage string, v ...interface{}) {
	if v != nil {
		log.Printf("[DEBUG]  "+ logMessage, v)
	} else {
		log.Printf("[DEBUG]  "+ logMessage)
	}
}

func Warn(logMessage string, v ...interface{}) {
	if v != nil {
		log.Printf("[WARN]  "+ logMessage, v)
	} else {
		log.Printf("[WARN]  "+ logMessage)
	} 
}

func Error(logMessage string, v ...interface{}) {
	if v != nil {
		log.Printf("[ERROR]  "+ logMessage, v)
	} else {
		log.Printf("[ERROR]  "+ logMessage)
	} 
}