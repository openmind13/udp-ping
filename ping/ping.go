package main

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

type UDPClient struct {
	conn *net.UDPConn
	mu   sync.RWMutex
	aims map[*net.UDPAddr]Aim
}

type Aim struct {
	Addr          *net.UDPAddr
	Stats         Stats
	PacketIndexes map[Index]struct{}
}

type Index struct {
	Uniq int
}

type Stats struct {
	AverageTime        time.Duration
	PacketsSendedTotal uint64
}

func NewUDPClient(remoteAddrs ...string) (*UDPClient, error) {
	c := &UDPClient{
		aims: make(map[*net.UDPAddr]Aim),
	}

	localAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 0,
	}

	for _, remoteAddrStr := range remoteAddrs {
		remoteAddr, err := net.ResolveUDPAddr("udp", remoteAddrStr)
		if err != nil {
			return nil, err
		}
		c.aims[remoteAddr] = Aim{
			Addr:          remoteAddr,
			Stats:         Stats{},
			PacketIndexes: make(map[Index]struct{}),
		}
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, err
	}
	c.conn = conn

	return c, nil
}

func (c *UDPClient) Start() {
	go func() {
		for {
			buffer := make([]byte, packet.MAX_SAFE_UDP_PAYLOAD_SIZE_BYTES)
			n, remoteAddr, err := c.conn.ReadFromUDP(buffer)
			if err != nil {
				log.Fatal(err)
			}

			incomePack := packet.Unmarshal(packet.ToPacketBin(buffer[:n]))
			if incomePack.Type != packet.ECHO_REPLY {
				log.Println("not echo reply packet")
				continue
			}

			c.mu.Lock()
			found := false
			for addr, aim := range c.aims {
				if remoteAddr.IP.Equal(addr.IP) && remoteAddr.Port == addr.Port {
					found = true
					delete(aim.PacketIndexes, Index{Uniq: incomePack.Index})
					log.Println("OK")
				}
			}
			c.mu.Unlock()

			if !found {
				log.Println("Addr not found:", remoteAddr)
				continue
			}
		}
	}()

	for {
		c.mu.RLock()
		for _, aim := range c.aims {
			index := rand.Int()
			echoRequest := packet.Packet{
				Index: index,
				Type:  packet.ECHO_REQUEST,
			}
			aim.PacketIndexes[Index{Uniq: index}] = struct{}{}
			aim.Stats.PacketsSendedTotal++
			_, err := c.conn.WriteToUDP(echoRequest.Marshal().ToSlice(), aim.Addr)
			if err != nil {
				log.Fatal(err)
			}
		}
		c.mu.RUnlock()

		time.Sleep(PING_PERIOD)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	log.Println("ping")

	client, err := NewUDPClient("127.0.0.1:9000")
	if err != nil {
		log.Fatal(err)
	}
	client.Start()
}
