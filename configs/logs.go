package configs

import "log"

type Logger struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

// func SetLogger(logger *Logger) {
// 	var LoggerInstance = &Logger{}
// 	LoggerInstance.ErrorLog = logger.ErrorLog
// 	LoggerInstance.InfoLog = logger.InfoLog
// }