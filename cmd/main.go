package main

import (
	"log"
	"math/rand"
	"time"
	"udp-ping/pinger"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	config := pinger.Config{
		CollectingStatsPeriod: 2 * time.Second,
		Aims: []pinger.AimConfig{
			{
				Name:       "Stanislav laptop",
				LocalAddr:  "0.0.0.0:3001",
				RemoteAddr: "192.168.8.184:3001",
				PingPeriod: 1 * time.Second,
			},
		},
	}
	if !config.IsValid() {
		log.Fatal("Pinger config not valid")
	}

	pinger, err := pinger.NewPinger(config)
	if err != nil {
		log.Fatal(err)
	}
	pinger.Start()
}
