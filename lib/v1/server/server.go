package server

import (
	"bufio"
	"fmt"
	"net"
)

type Server struct{}

// TODO Fix this to handle errors
func (s Server) Run() {
	ln, _ := net.Listen("tcp4", "0.0.0.0:2345")

	for {
		conn, _ := ln.Accept()
		go s.handleConnection(conn)
	}
}

func (s Server) handleConnection(conn net.Conn) {
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

type PacketType uint8

const (
	FILE_UPLOAD PacketType = iota << 1
)

type Packet struct {
	PacketType
	Header  any
	Payload any
}

func NewPacket(t PacketType) Packet {
	return Packet{
		PacketType: t,
		Header:     nil,
		Payload:    nil,
	}
}

// Protocol Header with:
// 1 - File size (8 bytes)
// 2 - Chunk size (4 bytes)
// 3 - Number of chunks (8 bytes)
// 4 - Filename length (2 bytes)
// 5 - Filename (variable)

type FileUploadHeader struct {
	FileSize    uint64 // 8 bytes
	ChunkSize   uint32 // 4 bytes
	ChunkCount  uint64 // 8 bytes
	FileNameLen uint16 // 2 bytes
	FileName    []byte // n bytes -> size from FileNameLen
}

// Chunks with:
// - Chunk ID (4 bytes)
// - Chunk CRC32 (32 bits => 4 bytes)
// - Chunk data (variable, as specified in header)
// With acknowledgments after each chunk

type FileChunkPayload struct {
	Id    uint32 // 4 bytes
	Crc32 uint32 // 4 bytes
	Data  []byte // n bytes
}
