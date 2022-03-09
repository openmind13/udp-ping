package pinger

import (
	"net"
)

type Statistics struct{}

type Pinger struct {
	Targets        map[string]*Target
	StatisticsChan chan Statistics
}

func NewPinger() (*Pinger, error) {
	pinger := Pinger{
		Targets:        map[string]*Target{},
		StatisticsChan: make(chan Statistics, 1),
	}

	t, err := NewTarget(nil, &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 9000,
	})
	if err != nil {
		return nil, err
	}
	pinger.Targets["0.0.0.0:9000"] = t
	go t.RunAsync()

	return &pinger, nil
}

func (p *Pinger) Start() {
	select {}
}
