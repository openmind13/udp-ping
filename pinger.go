package udppinger

import (
	"io"
	"time"
	"udp-ping/pinger/aim"

	"github.com/sirupsen/logrus"
)

type Statistics struct {
	AimStats map[string]aim.Stat
}

// func (s Statistics) Print() {
// 	logrus.WithFields(logrus.Fields{
// 		"aims_count": len(s.AimStats),
// 	}).Info("Statistics")
// 	for _, stat := range s.AimStats {
// 		logrus.WithFields(logrus.Fields{
// 			"avg_rtt":     stat.AverageRtt,
// 			"remote_addr": stat.RemoteAddr.String(),
// 			"rcv_pkt_cnt": stat.RecvPacketsCount,
// 		}).Info(stat.Name)
// 	}
// 	logrus.Info("-------------------------------------------------------------" +
// 		"-----------------------------------------")
// }

type Pinger struct {
	Config         Config
	Aims           map[string]*aim.Aim
	ConfigChan     chan Config
	StatisticsChan chan Statistics
	stopChan       chan struct{}

	logger io.Writer
}

func NewPinger(config Config, logger io.Writer) (*Pinger, error) {
	p := Pinger{
		Config:         config,
		Aims:           map[string]*aim.Aim{},
		ConfigChan:     make(chan Config, 1),
		StatisticsChan: make(chan Statistics, 1),
		stopChan:       make(chan struct{}, 1),
		logger:         logger,
	}

	for _, aimConfig := range config.Aims {
		if !aimConfig.IsValid() {
			logrus.Error("Aim config not valid ", aimConfig)
			continue
		}
		aim, err := aim.New(aimConfig, p.Config.PingInterval)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Failed create aim")
			continue
		}
		p.Aims[aim.Name] = aim
		go aim.Start()
	}

	return &p, nil
}

func (p *Pinger) Start() {
	statisticHook := time.NewTicker(p.Config.CollectingStatsInterval)

	for {
		select {
		case newConfig := <-p.ConfigChan:
			aimInConfig := map[string]bool{}
			for aimName := range p.Aims {
				aimInConfig[aimName] = false
			}
			for _, newAimConfig := range newConfig.Aims {
				aimInPinger, ok := p.Aims[newAimConfig.Name]
				if !ok {
					logrus.Debug("aim not found. create new aim ", newAimConfig.Name)
					aim, err := aim.New(newAimConfig, p.Config.PingInterval)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"error": err.Error(),
						}).Error("Failed create aim")
						continue
					}
					p.Aims[aim.Name] = aim
					go aim.Start()
				} else {
					aimInConfig[aimInPinger.Name] = true
					logrus.Debug("Aim exists ", aimInPinger.Name)
					aimInPinger.ConfigChan <- newAimConfig
				}
			}

			for aimName, flag := range aimInConfig {
				if !flag {
					aim := p.Aims[aimName]
					aim.Stop()
					delete(p.Aims, aimName)
				}
			}

		case <-statisticHook.C:
			statistics := Statistics{AimStats: make(map[string]aim.Stat)}
			for _, aim := range p.Aims {
				select {
				case stat := <-aim.StatChan:
					statistics.AimStats[aim.RemoteAddr.String()] = stat
				default:
					logrus.Warn("Statistic for aim [", aim.RemoteAddr.String(), "] not collected")
				}
			}
			p.StatisticsChan <- statistics

		case <-p.stopChan:
			for _, aim := range p.Aims {
				aim.Stop()
			}
			// close(p.stopChan)
			// close(p.ConfigChan)
			// close(p.StatisticsHookChan)
			close(p.StatisticsChan)
			logrus.Info("Pinger has stopped")
			return
		}

	}
}

func (p *Pinger) Stop() {
	p.stopChan <- struct{}{}
}
