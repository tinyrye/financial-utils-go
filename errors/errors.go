package errors

import (
	"log"
)

func PanicIf(err error, description string) {
	if err != nil {
		log.Panicf("Error during %s: %s", description, err)
	}
}
