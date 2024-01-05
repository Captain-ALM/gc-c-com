package transport

import (
	"log"
	"os"
)

func DebugPrintln(msg string) {
	if os.Getenv("DEBUG_COM") == "1" {
		log.Println("DEBUG_COM:", msg)
	}
}

func DebugErrIsNil(err error) bool {
	if err == nil {
		return true
	}
	DebugPrintln(err.Error())
	return false
}
