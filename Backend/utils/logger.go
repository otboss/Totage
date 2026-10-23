package utils

import (
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./app.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	})
}
