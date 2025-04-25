package main

import "github.com/tomascarruco/fileup/lib/v1/server"

func main() {
	server := server.Serve{}
	server.Run()
}
