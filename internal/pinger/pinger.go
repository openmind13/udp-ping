package pinger

import (
	"sync"
	"udp-ping/internal/config"
)

type Pinger struct {
	mu     sync.RWMutex
	config config.Config
}

func New(cfg config.Config) (*Pinger, error) {
	p := &Pinger{}

	return p, nil
}

func (p *Pinger) PingAims() error {
	return nil
}

func (p *Pinger) startUdpServer() error {
	return nil
}

func (p *Pinger) GetConfig() config.Config {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.config
}

func (p *Pinger) ChangeConfig(cfg config.Config) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.config = cfg
	// TODO
}
