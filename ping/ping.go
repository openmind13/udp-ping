package main

import (
	"log"
	"net"
	"time"
)

var (
	port = 9000
)

type Client struct {
	conn *net.UDPConn
}

func main() {
	log.Println("ping")

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		Port: port,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("Connected")

	for {
		msg := []byte("test")
		// _, err := conn.WriteToUDP(msg, &net.UDPAddr{
		// 	Port: port,
		// })

		_, err := conn.Write(msg)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("write")
		time.Sleep(time.Second)
	}
}
