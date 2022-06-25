package config

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

var (
	ConfigChan = make(chan Config, 1)
)

type Config struct {
	LogLevel   string        `mapstructure:"log_level"`
	PingPeriod time.Duration `mapstructure:"ping_period" validate:"required"`
	Aims       []Aim         `mapstructure:"aim" json:"aims" yaml:"aim"`
}

type Aim struct {
	Name       string `mapstructure:"name" json:"name" yaml:"name" validate:"required"`
	RemoteAddr string `mapstructure:"remote_addr" json:"remote_addr" yaml:"remote_addr" validate:"required"`
	LocalAddr  string `mapstructure:"local_addr" json:"local_addr" yaml:"local_addr"`
}
