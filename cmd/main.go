package main

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"udp-ping/configurator"
	"udp-ping/logger"
	"udp-ping/pinger"

	"github.com/sirupsen/logrus"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	logger.Init()

	if err := configurator.Start(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal(err)
	}

	config := <-configurator.ConfigChan
	if !config.IsValid() {
		logrus.Fatal("config not valid")
	}

	logger.LogLevelChan <- config.LogLevel

	logrus.Info("Starting pinger")

	myPinger, err := pinger.NewPinger(config)
	if err != nil {
		logrus.Fatal(err)
	}

	go myPinger.Start()

	go func() {
		for {
			stat, ok := <-myPinger.StatisticsChan
			if !ok {
				logrus.Info("Pinger stat channel closed")
				return
			}
			stat.Print()
		}
	}()

	go func() {
		time.Sleep(time.Second * 10)
		myPinger.Stop()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case newConfig := <-configurator.ConfigChan:
			logger.LogLevelChan <- newConfig.LogLevel
			myPinger.ConfigChan <- newConfig

		case sig := <-sigChan:
			logrus.Info("Stopped due the signal: " + sig.String())
			os.Exit(0)
		}
	}

}
