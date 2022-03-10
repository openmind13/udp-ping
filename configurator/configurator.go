package configurator

import (
	"udp-ping/logger"
	"udp-ping/pinger"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	ConfigChan = make(chan pinger.Config, 1)
)

func Start() error {
	viper.SetConfigName("config.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	var config pinger.Config
	if err := viper.Unmarshal(&config); err != nil {
		logrus.Fatal(err)
	}

	ConfigChan <- config

	viper.WatchConfig()

	viper.OnConfigChange(func(in fsnotify.Event) {
		logrus.Debug("Changes in config!!!")
		newConfig := pinger.Config{}
		if err := viper.Unmarshal(&newConfig); err != nil {
			logrus.Fatal(err)
		}
		logger.LogLevelChan <- config.LogLevel
		ConfigChan <- newConfig
	})

	return nil
}
