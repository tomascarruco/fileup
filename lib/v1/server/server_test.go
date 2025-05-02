package server

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/tomascarruco/fileup/lib/v1/protocol"
)

const (
	DEFAULT_SERVER_ADDRESS = "0.0.0.0:2345"
)

func TestMain(m *testing.M) {
	fmt.Println("Setting up...")

	server := Server{}
	{
		go server.Run()
		time.Sleep(100 * time.Millisecond)
	}

	exitCode := m.Run()

	fmt.Println("Exiting...")

	os.Exit(exitCode)
}

func TestHandlePingRequest(t *testing.T) {
	conn, err := net.Dial("tcp4", DEFAULT_SERVER_ADDRESS)
	if err != nil {
		t.Fatal("Failed to dial to expected ip and port\n")
	}

	packet := protocol.NewPacket(protocol.PING)

	// Encode new PING packet
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(packet)
	if err != nil {
		t.Fatalf("Failed to encode packet: %s\n", err.Error())
	}

	// Decode incoming packet response
	decoder := gob.NewDecoder(conn)
	var packetRec protocol.Packet
	err = decoder.Decode(&packetRec)
	if err != nil {
		t.Fatalf("Failed to decode packet\n")
	}

	if packetRec.PacketType != protocol.PONG {
		t.Fatalf("Invalid packet returned, should be PONG, was: %d\n", packet.PacketType)
	}
}
