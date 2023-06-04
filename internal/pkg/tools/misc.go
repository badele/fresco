package tools

import "log"

func CheckError(err error) {
	if err != nil {
		log.Printf("error: %v", err)
	}
}
