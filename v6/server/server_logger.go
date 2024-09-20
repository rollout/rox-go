package server

import (
	"log"

	"github.com/rollout/rox-go/v6/core/logging"
)

type serverLogger struct {
}

func NewServerLogger() logging.Logger {
	return &serverLogger{}
}

func (*serverLogger) Debug(message string, err interface{}) {
	if err != nil {
		log.Println("DEBUG: ", message, err)
	} else {
		log.Println("DEBUG: ", message)
	}
}

func (*serverLogger) Warn(message string, err interface{}) {
	if err != nil {
		log.Println("WARN: ", message, err)
	} else {
		log.Println("WARN: ", message)
	}
}

func (*serverLogger) Error(message string, err interface{}) {
	if err != nil {
		log.Println("ERROR: ", message, err)
	} else {
		log.Println("ERROR: ", message)
	}
}
