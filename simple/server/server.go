package main

import (
	"flag"
	"log"
	"net"
	"udp-ping/packet"
)

var (
	localAddrFlag = flag.String("addr", "0.0.0.0:9000", "local udp addr")
)

func main() {
	flag.Parse()

	localAddr, err := net.ResolveUDPAddr("udp", *localAddrFlag)
	if err != nil {
		log.Fatal(err)
	}

	udpServer, err := NewServer(localAddr)
	if err != nil {
		log.Fatal(err)
	}
	udpServer.Start()
}

type UDPServer struct {
	conn *net.UDPConn
}

func NewServer(addr *net.UDPAddr) (*UDPServer, error) {
	s := &UDPServer{}

	log.Println("Starting listening udp on:", addr)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	s.conn = conn

	return s, nil
}

func (s *UDPServer) Start() {
	for {
		buffer := make([]byte, packet.MAX_SAFE_UDP_PAYLOAD_SIZE_BYTES)

		n, remoteAddr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}

		incomingPacket := packet.Unmarshal(packet.ToPacketBin(buffer[:n]))
		if incomingPacket.Type != packet.ECHO_REQUEST {
			log.Println("not echo request")
			continue
		}

		echoReplyPacket := packet.Packet{
			Index: incomingPacket.Index,
			Time:  incomingPacket.Time,
			Type:  packet.ECHO_REPLY,
		}
		_, err = s.conn.WriteToUDP(echoReplyPacket.Marshal().ToSlice(), remoteAddr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
