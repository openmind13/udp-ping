package pinger

import "time"

type Config struct {
	CollectingStatsPeriod time.Duration
	Aims                  []AimConfig
}

func (c Config) IsValid() bool {
	if c.CollectingStatsPeriod == 0 || c.Aims == nil {
		return false
	}
	return true
}

type AimConfig struct {
	Name       string
	LocalAddr  string
	RemoteAddr string
	PingPeriod time.Duration
}

func (ac AimConfig) IsValid() bool {
	if ac.Name == "" || ac.LocalAddr == "" || ac.RemoteAddr == "" || ac.PingPeriod == 0 {
		return false
	}
	return true
}
