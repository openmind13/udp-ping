package aim

import (
	"net"
	"sync"
	"udp-ping/internal/config"
)

type Aim struct {
	mu     sync.RWMutex
	config config.Aim

	localAddr  net.UDPAddr
	remoteAddr net.UDPAddr
	conn       *net.UDPConn
}

func New(cfg config.Aim) (*Aim, error) {
	a := &Aim{
		config: cfg,
	}

	localAddr, err := net.ResolveUDPAddr("udp", a.config.LocalAddr)
	if err != nil {
		return nil, err
	}
	a.localAddr = *localAddr

	return a, nil
}

func (a *Aim) Start() {}
