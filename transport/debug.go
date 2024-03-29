package transport

import (
	"log"
	"os"
)

func debugPrintln(msg string) {
	if os.Getenv("DEBUG_COM") == "1" {
		log.Println("DEBUG_COM:", msg)
	}
}

func debugErrIsNil(err error) bool {
	if err == nil {
		return true
	}
	debugPrintln(err.Error())
	return false
}
