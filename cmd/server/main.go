package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/netip"
)

func main() {
	initFlag()

	ipAddr, err := netip.ParseAddrPort(ServerFlags.Ip + ":" + ServerFlags.Port)
	if err != nil {
		log.Fatal("Failed to parse the IP and Port")
	}
	fmt.Printf("Running on: %s:%d\n", ipAddr.Addr(), ipAddr.Port())

	ln, err := net.Listen("tcp4", ipAddr.String())
	if err != nil {
		log.Fatal("Failed to parse the IP and Port")
	}

	for {
		fmt.Println("Accepting...")
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Could not accept connection!!!")
			continue
		}
		fmt.Println("Accepted!!")
		go handleConnection(conn)
	}
}

func handleConnection(con net.Conn) {
	reader := bufio.NewReader(con)

	content, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Failed to read the content from the connection")
	}

	fmt.Printf("Content: %s\n", content)
	_ = con.Close()
}

type serverFlags struct {
	Port      string
	Ip        string
	IpVersion string
}

var ServerFlags = serverFlags{}

func initFlag() {
	const (
		defaultPort = "8080"
		defaultIp   = "0.0.0.0"
		defaultIPv  = "v4"
	)

	flag.StringVar(&ServerFlags.Port, "port", defaultPort, "server listening port")
	flag.StringVar(&ServerFlags.Ip, "ip", defaultIp, "server listening IP")
	flag.StringVar(&ServerFlags.IpVersion, "ipv", defaultIPv, "IPv 4 or 6")

	flag.Parse()
}
