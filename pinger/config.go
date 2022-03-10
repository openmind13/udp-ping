package pinger

import (
	"time"
	"udp-ping/pinger/aim"
)

type Config struct {
	LogLevel              string          `mapstructure:"log_level"`
	CollectingStatsPeriod time.Duration   `mapstructure:"collecting_stats_period"`
	PingPeriod            time.Duration   `mapstructure:"ping_period"`
	Aims                  []aim.AimConfig `mapstructure:"aim"`
}

func (c Config) IsValid() bool {
	if c.CollectingStatsPeriod == 0 || c.Aims == nil || c.PingPeriod == 0 {
		return false
	}
	return true
}
