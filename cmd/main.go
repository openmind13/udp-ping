package main

import (
	"log"
	"math/rand"
	"time"
	"udp-ping/pinger"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	pinger, err := pinger.NewPinger()
	if err != nil {
		log.Fatal(err)
	}
	pinger.Start()
}
