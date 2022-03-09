package packet

import (
	"bytes"
	"encoding/gob"
	"time"
)

const (
	MAX_SAFE_UDP_PAYLOAD_SIZE_BYTES = 508

	PACKET_TIMESTAMP_SIZE_BYTES = 8
	PACKET_TYPE_SIZE_BYTES      = 4
	PACKET_PAYLOAD_SIZE_BYTES   = 496
)

type PacketBin [MAX_SAFE_UDP_PAYLOAD_SIZE_BYTES]byte

func (pb PacketBin) ToSlice() []byte {
	return pb[:]
}

func ToPacketBin(buffer []byte) PacketBin {
	var data [MAX_SAFE_UDP_PAYLOAD_SIZE_BYTES]byte
	copy(data[:], buffer)
	return data
}

type Packet struct {
	Index int
	Time  time.Time
	Type  PacketType
}

type PacketType int

const (
	ECHO_REQUEST = 0
	ECHO_REPLY   = 1
)

type PacketData struct{}

func (p Packet) Marshal() PacketBin {
	var bin PacketBin

	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(p); err != nil {
		return PacketBin{}
	}
	copy(bin[:], buffer.Bytes())
	return bin
}

func Unmarshal(binData PacketBin) Packet {
	p := Packet{}
	var buffer bytes.Buffer
	buffer.Write(binData[:])
	if err := gob.NewDecoder(&buffer).Decode(&p); err != nil {
		return Packet{}
	}
	return p
}
