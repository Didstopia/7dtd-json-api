package main

import (
	"log"

	"github.com/reiver/go-telnet"
)

func main() {
	var caller telnet.Caller = telnet.StandardCaller

	//@TOOD: replace "example.net:5555" with address you want to connect to.
	log.Println("Connecting to telnet server..")
	telnet.DialToAndCall("example.net:5555", caller)
}
