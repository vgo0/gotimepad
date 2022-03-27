package handler

import "log"

func CheckError(err error) {
	if err != nil {
		log.Fatalf("%s\n", err)
	}
}