package server

import (
	"encoding/gob"
	"log"
	"net"
	"reflect"

	"github.com/tomascarruco/fileup/lib/v1/protocol"
)

/* TODO
   Setup packet parssing to handle diferent types of actions
*/

type Serve struct {
}

func (s Serve) Run() {
	// TODO Fix this to handle errors
	ln, _ := net.Listen("tcp4", "0.0.0.0:2345")

	for {
		conn, _ := ln.Accept()
		go s.processNewConnection(conn)
	}
}

func (s Serve) processNewConnection(conn net.Conn) {
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
				"Skipping this packet type: %s",
				reflect.TypeOf(packetType).Name(),
			)
			continue

		case protocol.PING:
			packet = protocol.NewPacket(protocol.PONG_RESPONSE)

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
