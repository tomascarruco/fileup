package server

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"reflect"

	"github.com/tomascarruco/fileup/lib/v1/protocol"
)

/* TODO
   Setup packet parssing to handle diferent types of actions
*/

const (
	DEFAULT_SERVER_PORT      = 8765
	DEFAULT_SERVER_LISTEN_IP = "0.0.0.0"
)

type Server struct {
	ServerConfigurable

	listener        net.Listener
	poolWorkerCount uint
}

type ServerOption func(*Server)

func NewServer(opt ...ServerOption) *Server {
	c := new(Server)

	c.Network = struct {
		Port uint16
		IP   string
	}{
		Port: DEFAULT_SERVER_PORT,
		IP:   DEFAULT_SERVER_LISTEN_IP,
	}

	c.Computing.MaxWorkers = 255

	for _, o := range opt {
		o(c)
	}

	c.poolWorkerCount = 0

	if c.listener != nil {
		return c
	}

	baseAddress := fmt.Sprintf("%s:%d", c.Network.IP, c.Network.Port)

	address := net.ParseIP(baseAddress)
	if address == nil {
		address = net.ParseIP("0.0.0.0:8765")
	}

	listener, err := net.Listen("tcp4", address.String())
	if err != nil {
		log.Printf("Could not bind to address: %s", err.Error())
	}
	c.listener = listener

	return c
}

func (s Server) Run() {
	// TODO Fix this to handle errors
	for {
		conn, _ := s.listener.Accept()
		go s.processNewConnection(conn)
	}
}

func (s Server) processNewConnection(conn net.Conn) {
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
