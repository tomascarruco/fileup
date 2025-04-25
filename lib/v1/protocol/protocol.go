package protocol

import "errors"

type PacketType uint8

const (
	// File management related packet types
	FILE_UP_REQUEST PacketType = iota + 1
	FILE_UP_CHUNK
	FOLDER_CREATE
	FOLDER_DELETE
	// End ---
	// Status related packet types
	Ping
	Pong
)

// Pakcets are sent in Littele Endian order
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

func NewFileUploadHeader(fileSize uint64, chunkSize uint32, fileName string) (FileUploadHeader, error) {
	if len(fileName) > 255 {
		return FileUploadHeader{}, errors.New("File name to big")
	}

	return FileUploadHeader{
		FileSize:    fileSize,
		ChunkSize:   chunkSize,
		ChunkCount:  fileSize / uint64(chunkSize),
		FileNameLen: uint16(len(fileName)),
		FileName:    []byte(fileName),
	}, nil
}

// Chunks with:
// - Chunk ID (4 bytes)
// - Chunk CRC32 (32 bits => 4 bytes)
// - Chunk data (variable, as specified in header)
// With acknowledgments after each chunk
//
// Probably a pool is a good place to store this kind of object

type FileChunkPayload struct {
	Id    uint32 // 4 bytes
	Crc32 uint32 // 4 bytes
	Data  []byte // n bytes
}
