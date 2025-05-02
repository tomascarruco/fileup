package protocol

import (
	"bufio"
	"errors"
	"hash/crc32"
	"log"
	"path/filepath"
)

/*
Ping and Pong requests
Create Folders (with permissions?)
Create Files, then start sync process:
	1. Is the file new? If so upload it (store each upload step to allow resume)
	2. After file upload check its integrity
	3. Finaly, keep the file stored
*/

const (
	FILE_NAME_MAX_LEN   int = 255
	FILE_NAME_MIN_LEN   int = 1
	FOLDER_PATH_MAX_LEN int = 4096
	FOLDER_PATH_MIN_LEN int = 1
)

type PacketType uint8

const (
	// File management related packet types
	FILE_UP_DEFINE PacketType = iota + 1
	FILE_UP_CHUNK
	// Folder management related packet types
	FOLDER_CREATE
	FOLDER_DELETE
	// End ---
	// Status related packet types
	PING
	PONG
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

// Folder creation:
// 1 - Folder Name Len (16 bytes)
// 2 - Folder Name (variable)
// 3 - Folder Path Len (16 bytes)
// 4 - Folder Path (variable)

type FolderCreateInfo struct {
	FolderPathLength uint16
	FolderPath       []byte
}

func NewFolderCreateInfo(path string) (FolderCreateInfo, error) {

	folderPathLength := len(path)

	switch {
	case folderPathLength > FOLDER_PATH_MAX_LEN:
		return FolderCreateInfo{},
			ProtocolError{What: FolderPathToLargeError}

	case folderPathLength < FOLDER_PATH_MIN_LEN:
		return FolderCreateInfo{},
			ProtocolError{What: FolderPathToSmallError}
	}

	cleanPath := filepath.Clean(path)

	if cleanPath == "." {
		return FolderCreateInfo{},
			ProtocolError{What: LocalFolderAsPathError}
	}

	return FolderCreateInfo{
		FolderPathLength: uint16(folderPathLength),
		FolderPath:       []byte(path),
	}, nil
}

// Protocol Header with:
// 1 - File size (8 bytes)
// 2 - Chunk size (4 bytes)
// 3 - Number of chunks (8 bytes)
// 4 - Filename length (2 bytes)
// 5 - Filename (variable)

type FileUploadInfo struct {
	FileSize    uint64 // 8 bytes
	ChunkSize   uint64 // 8 bytes
	ChunkCount  uint64 // 8 bytes
	FileNameLen uint16 // 2 bytes
	FileName    []byte // n bytes -> size from FileNameLen
}

const (
	CHUNK_SIZE_SMALL  uint64 = 1 * 1024 * 1024       // 1MB
	CHUNK_SIZE_MEDIUM        = 5 * CHUNK_SIZE_SMALL  // 5MB
	CHUNK_SIZE_LARGE         = 10 * CHUNK_SIZE_SMALL // 10MB
)

func NewFileUploadInfo(
	fileSize uint64,
	chunkSize uint64,
	fileName string,
) (FileUploadInfo, error) {

	fileNameLength := len(fileName)

	switch {
	case fileNameLength > FILE_NAME_MAX_LEN:
		return FileUploadInfo{},
			ProtocolError{
				What: FileNameToLargeError,
				How:  "File name is larger than 255 bytes, please reduce it",
			}

	case fileNameLength < FILE_NAME_MIN_LEN:
		return FileUploadInfo{},
			ProtocolError{
				What: FileNameToSmallError,
				How:  "File name cant be empty",
			}

	case chunkSize > fileSize:
		return FileUploadInfo{},
			ProtocolError{
				What: ChunkLargerThanFileError,
				How:  "A single chunk size can't be larger than the size of the file",
			}
	}

	return FileUploadInfo{
		FileSize:    fileSize,
		ChunkSize:   chunkSize,
		ChunkCount:  fileSize / uint64(chunkSize),
		FileNameLen: uint16(fileNameLength),
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

type FileChunkData struct {
	Id    uint32 // 4 bytes
	Crc32 uint32 // 4 bytes
	Data  []byte // n bytes
}

func NewFileChunk(id uint32, data bufio.Reader, chunkSize uint64) (FileChunkData, error) {
	chunk := FileChunkData{
		Id:   id,
		Data: make([]byte, chunkSize, 0x00),
	}

	// Since the buffer is created with make,
	// bufio.Reader.Read, will ever overflow the buffer
	// Other operations might;
	_, err := data.Read(chunk.Data)
	if errors.Is(err, bufio.ErrBufferFull) {
		log.Printf("failed to read chunck into buffer, should not happen")

		return FileChunkData{},
			ProtocolError{
				What: ChunkLargerThanFileError,
				How:  "input read larger than allocated chunk",
			}
	}

	chunk.Crc32 = crc32.Checksum(chunk.Data, crc32.IEEETable)

	return chunk, nil
}
