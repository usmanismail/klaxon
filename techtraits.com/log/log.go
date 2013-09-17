package log

import (
	"appengine"
	"log"
)

func Info(logMessage string, v ...interface{}) {
	log.Printf("[INFO]  "+logMessage, v)
}

func Debug(logMessage string, v ...interface{}) {
	log.Printf("[DEBUG]  "+logMessage, v)
}

func Warn(logMessage string, v ...interface{}) {
	log.Printf("[WARN]  "+logMessage, v)
}

func Error(logMessage string, v ...interface{}) {
	log.Printf("[ERROR]  "+logMessage, v)
}

func Infof(context appengine.Context, logMessage string, v ...interface{}) {
	context.Infof("[INFO]  "+logMessage, v)
}

func Debugf(context appengine.Context, logMessage string, v ...interface{}) {
	context.Debugf("[DEBUG]  "+logMessage, v)
}

func Warnf(context appengine.Context, logMessage string, v ...interface{}) {
	context.Warningf("[WARN]  "+logMessage, v)
}

func Errorf(context appengine.Context, logMessage string, v ...interface{}) {
	context.Errorf("[ERROR]  "+logMessage, v)
}

func Criticalf(context appengine.Context, logMessage string, v ...interface{}) {
	context.Criticalf("[ERROR]  "+logMessage, v)
}
