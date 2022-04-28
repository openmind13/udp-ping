package udppinger

import (
	"time"
	"udp-ping/aim"
)

type Config struct {
	CollectingStatsInterval time.Duration   `mapstructure:"collecting_stats_interval" json:"collecting_stats_interval"`
	PingInterval            time.Duration   `mapstructure:"ping_interval" json:"ping_inteval"`
	Aims                    []aim.AimConfig `mapstructure:"aim"`
}

func (c Config) IsValid() bool {
	if c.CollectingStatsInterval == 0 || c.Aims == nil || c.PingInterval == 0 {
		return false
	}
	return true
}
