package pinger

import (
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
	"udp-ping/packet"
)

const (
	PING_PERIOD = time.Second
)

type Target struct {
	Mu         sync.RWMutex
	Conn       *net.UDPConn
	RemoteAddr *net.UDPAddr
	Stats      Stats

	SendPacket chan packet.Packet
	RecvPacket chan packet.Packet

	PacketIndexes map[Index]time.Time
}

type Index struct {
	int
}

type Stats struct {
	SentPacketsTotal int
	RecvPacketsTotal int
	AverageTime      time.Duration
	StartTime        time.Time
}

func NewTarget(localAddr *net.UDPAddr, remoteAddr *net.UDPAddr) (*Target, error) {
	t := &Target{
		PacketIndexes: make(map[Index]time.Time),

		SendPacket: make(chan packet.Packet, 1),
		RecvPacket: make(chan packet.Packet, 1),
	}
	t.RemoteAddr = remoteAddr

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, err
	}
	t.Conn = conn

	go func() {
		for {
			select {
			case pack := <-t.SendPacket:
				t.PacketIndexes[Index{pack.Index}] = pack.Time
				t.Stats.SentPacketsTotal++

			case pack := <-t.RecvPacket:
				sendTime, ok := t.PacketIndexes[Index{pack.Index}]
				if !ok {
					log.Println("packet index not found", pack.Index)
				} else {
					t.Stats.RecvPacketsTotal++
					packetTripTime := time.Since(sendTime)
					t.Stats.AverageTime = time.Since(t.Stats.StartTime) / time.Duration(t.Stats.RecvPacketsTotal)
					log.Println("OK", packetTripTime, t.Stats.AverageTime)
				}
			}
		}
	}()

	return t, nil
}

func (t *Target) RunAsync() {
	go func() {
		for {
			buffer := make([]byte, packet.MAX_SAFE_UDP_PAYLOAD_SIZE_BYTES)
			n, _, err := t.Conn.ReadFromUDP(buffer)
			if err != nil {
				log.Fatal(err)
			}
			incomePacket := packet.Unmarshal(packet.ToPacketBin(buffer[:n]))
			if incomePacket.Type == packet.ECHO_REPLY {
				t.RecvPacket <- incomePacket
			} else {
				log.Println("not echo packet")
			}
		}
	}()

	t.Stats.StartTime = time.Now()
	for {
		index := rand.Int()
		echoPacket := packet.Packet{
			Index: index,
			Type:  packet.ECHO_REQUEST,
			Time:  time.Now(),
		}
		t.SendPacket <- echoPacket
		_, err := t.Conn.WriteToUDP(echoPacket.Marshal().ToSlice(), t.RemoteAddr)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(PING_PERIOD)
	}
}

// func (t *Target) Run() {
// 	for {
// 		index := rand.Int()
// 		echoRequest := packet.Packet{
// 			Index: index,
// 			Type:  packet.ECHO_REQUEST,
// 		}

// 		t.PacketIndexes[Index{index}] = time.Now()
// 		_, err := t.Conn.WriteToUDP(echoRequest.Marshal().ToSlice(), t.RemoteAddr)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		buffer := make([]byte, packet.MAX_SAFE_UDP_PAYLOAD_SIZE_BYTES)
// 		n, _, err := t.Conn.ReadFromUDP(buffer)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		incomePacket := packet.Unmarshal(packet.ToPacketBin(buffer[:n]))
// 		if incomePacket.Type != packet.ECHO_REPLY {
// 			log.Println("not echo packet")
// 			continue
// 		}

// 		sendTime, ok := t.PacketIndexes[Index{incomePacket.Index}]
// 		if !ok {
// 			log.Println("packet index not found", incomePacket.Index)
// 			continue
// 		}
// 		log.Println("OK", time.Since(sendTime))

// 		time.Sleep(PING_PERIOD)
// 	}
// }
