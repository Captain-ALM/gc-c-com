package packet

import (
	"log"
	"os"
)

func debugPrintln(msg string) {
	if os.Getenv("DEBUG_COM_PK") == "1" {
		log.Println("DEBUG_COM_PK:", msg)
	}
}

func debugErrIsNil(err error) bool {
	if err == nil {
		return true
	}
	debugPrintln(err.Error())
	return false
}
