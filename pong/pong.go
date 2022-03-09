package main

import (
	"log"
	"net"
	"sync"
	"time"
	"udp-ping/packet"
)

var (
	localAddr = net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 9000,
	}
)

type Client struct {
	Addr        *net.UDPAddr
	AverageTime time.Duration
	PacketCount int
}

type UDPServer struct {
	conn    *net.UDPConn
	mu      sync.RWMutex
	clients map[string]Client
}

func NewServer(addr net.UDPAddr) (*UDPServer, error) {
	s := &UDPServer{
		clients: make(map[string]Client),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return nil, err
	}
	s.conn = conn

	return s, nil
}

func (s *UDPServer) Start() {
	for {
		buffer := make([]byte, packet.MAX_SAFE_UDP_PAYLOAD_SIZE_BYTES)

		_, remoteAddr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}
		s.mu.Lock()
		if _, ok := s.clients[remoteAddr.String()]; !ok {
			s.clients[remoteAddr.String()] = Client{
				Addr: remoteAddr,
			}
		}
		s.mu.Unlock()

		pack := packet.Unmarshal(packet.ToPacketBin(buffer))
		if pack.Type != packet.ECHO_REQUEST {
			log.Println("not echo request")
			continue
		}

		pongMsg := packet.Packet{
			Index: pack.Index,
			Type:  packet.ECHO_REPLY,
		}
		_, err = s.conn.WriteToUDP(pongMsg.Marshal().ToSlice(), remoteAddr)
		if err != nil {
			log.Fatal(err)
		}

		// log.Println(remoteAddr, buffer[:n], time.Since(readTime))
	}
}

func main() {
	log.Println("pong")

	udpServer, err := NewServer(localAddr)
	if err != nil {
		log.Fatal(err)
	}
	udpServer.Start()
}
