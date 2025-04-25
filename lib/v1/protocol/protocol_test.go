package protocol

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestPacketEncodingDecoding(t *testing.T) {
	packet := Packet{
		PacketType: FILE_UP_REQUEST,
		Header: FileUploadHeader{
			FileSize:    8 * 4024,
			ChunkSize:   1024,
			ChunkCount:  (8 * 4024) / 1024,
			FileNameLen: 10,
			FileName:    []byte("0123456789"),
		},
	}

	var network bytes.Buffer

	gob.Register(packet.Header)

	if packet.Payload != nil {
		gob.Register(packet.Payload)
	}
	enc := gob.NewEncoder(&network)

	err := enc.Encode(packet)
	if err != nil {
		t.Fatalf("Failed to encode packet, error: %s", err.Error())
	}

	dec := gob.NewDecoder(&network)
	var packetReceived Packet
	err = dec.Decode(&packetReceived)
	if err != nil {
		t.Fatalf("Failed to decode packet, error: %s", err.Error())
	}

	originalHeader := packet.Header.(FileUploadHeader)
	receivedHeader := packetReceived.Header.(FileUploadHeader)

	if originalHeader.FileNameLen != receivedHeader.FileNameLen {
		t.Fail()
	}
}
