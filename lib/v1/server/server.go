package server

import (
	"bufio"
	"fmt"
	"net"
)

type Serve struct{}

// TODO Fix this to handle errors
func (s Serve) Run() {
	ln, _ := net.Listen("tcp4", "0.0.0.0:2345")

	for {
		conn, _ := ln.Accept()
		go s.handleConnection(conn)
	}
}

func (s Serve) handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		text, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			break
		}
		fmt.Printf("-> %s", text)
	}
}
