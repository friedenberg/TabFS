package main

import (
	"flag"
	"log"
	"os/exec"
)

func main() {
	log.SetPrefix("chrome-json ")

	flagAddr := flag.String("addr", "chrome.sock", "address to serve from")

	log.Printf("starting server on %s", *flagAddr)

	socket := ServerSocket{SockPath: *flagAddr}
	socket.Listen()

	socat := exec.Command(
		"socat",
		"TCP4-LISTEN:3001,fork,bind=127.0.0.1",
		"UNIX-CONNECT:chrome.sock",
	)

	socat.Start()

	server := Server{}
	server.Serve(socket.Listener)
}
