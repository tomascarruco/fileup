package lib

import (
	"math/rand"
	"net/netip"
)

type RequestStatus uint8

const (
	Synced = iota
	UnSynced
	Syncing
)

type Request struct {
	Id     uint8
	Port   uint16
	Status RequestStatus
}

type FileServer struct {
	Address  netip.Addr
	Port     uint16
	requests map[uint8]Request
}

func NewFileServer(address netip.AddrPort) (fs FileServer) {
	fs.Address = address.Addr()
	fs.Port = address.Port()
	fs.requests = make(map[uint8]Request, 10)

	return
}

func (fs *FileServer) NewRequest(port uint16) uint8 {
	return uint8(rand.Uint32())
}
