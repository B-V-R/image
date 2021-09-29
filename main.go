package main

import (
	"image/client"
	"image/server"
	"os"
)

func main() {
	args := os.Args

	if len(args) > 1 {
		client := client.NewClient("0.0.0.0:8080", "certs/localhost.cert")
		client.Upload(args[1])
	} else {
		cert := "certs/localhost.cert"
		key := "certs/localhost.key"

		server.StartServer("0.0.0.0:8080", "img", cert, key)
	}
}
