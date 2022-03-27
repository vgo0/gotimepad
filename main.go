package main

import (
	"log"
	"os"

	"github.com/vgo0/gotimepad/create"
	"github.com/vgo0/gotimepad/decrypt"
	"github.com/vgo0/gotimepad/encrypt"
)

func main() {
	err_string := "Expected 'init (i)', 'encode (e)', 'decode (d)' command"

	if len(os.Args) < 2 {
		log.Fatalln(err_string)
	}

	switch os.Args[1] {
	case "init":
	case "i":
		create.Exec(os.Args[2:])
	case "encode":
	case "e":
		encrypt.Execute(os.Args[2:])
	case "decode":
	case "d":
		decrypt.Execute(os.Args[2:])
	default:
		log.Fatalln(err_string)
	}
	
}