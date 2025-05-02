package server

import (
	"encoding/gob"
	"log"
	"net"
	"reflect"
	"sync"

	"github.com/tomascarruco/fileup/lib/v1/protocol"
)

/* TODO
   Setup packet parssing to handle diferent types of actions
*/

type Server struct {
	listener *net.Listener
	ConnInfo struct {
		Port uint16
		IP   string
	}
	wg              sync.WaitGroup
	poolWorkerCount uint
	MaxPools        uint
	MinPools        uint
}

type ServerOption func(*Server)

func NewServer(opt ...ServerOption) *Server {
	c := new(Server)

	c.MaxPools = 255
	c.MinPools = 1

	for _, o := range opt {
		o(c)
	}

	c.poolWorkerCount = 0

	if c.listener != nil {
		// c.listener = net.ListenIP()
	}

	return c
}

func (s *Server) Run() {
	// TODO Fix this to handle errors
	ln, _ := net.Listen("tcp4", "0.0.0.0:2345")

	for {
		conn, _ := ln.Accept()
		go s.processNewConnection(conn)
	}
}

func (s *Server) processNewConnection(conn net.Conn) {
	// Proabably best put in a pool
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	for {
		// Receive and decode incoming packet
		var packetReceived protocol.Packet
		err := decoder.Decode(&packetReceived)
		if err != nil {
			log.Fatalf("Failed to decode received packet: %s\n", err.Error())
		}

		packetType := packetReceived.PacketType

		var packet protocol.Packet

		// Handle packet type
		switch packetType {
		default:
			log.Printf(
				"(UNIMPLEMENTED) Skipping this packet type: %s",
				reflect.TypeOf(packetType).Name(),
			)
			continue

		case protocol.PING:
			packet = protocol.NewPacket(protocol.PONG)
		}

		switch {
		case packet.Header != nil:
			gob.Register(packet.Header)

		case packet.Payload != nil:
			gob.Register(packet.Payload)
		}

		err = encoder.Encode(packet)
		if err != nil {
			log.Fatalf("Failed to encode packet: %s\n", err.Error())
		}
	}
}
