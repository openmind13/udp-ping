package main

import (
	"log"
	"net"
)

var (
	port = 9000
)

type Server struct {
	listener *net.UDPConn
}

func main() {
	log.Println("pong")

	log.Println("listen udp: ", port)
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		buffer := []byte{}
		buffer2 := []byte{}

		a, b, c, addr, err := conn.ReadMsgUDP(buffer, buffer2)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(a, b, c, addr)
		log.Println(buffer, buffer2)
	}
}
