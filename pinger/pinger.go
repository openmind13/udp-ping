package pinger

import (
	"fmt"
	"log"
	"time"
)

type Statistics struct {
	AimStats map[string]AimStat
}

func (s Statistics) Print() {
	fmt.Println()
	for _, stat := range s.AimStats {
		fmt.Println(stat.Name, stat.RemoteAddr.String(), stat.AverageRtt)
	}
}

type Pinger struct {
	CollectingStatsPeriod time.Duration
	Aims                  map[string]*Aim
	ConfigChan            chan<- Config
	StatisticsChan        <-chan Statistics
}

func NewPinger(config Config) (*Pinger, error) {
	pinger := Pinger{
		CollectingStatsPeriod: config.CollectingStatsPeriod,
		Aims:                  map[string]*Aim{},
		StatisticsChan:        make(chan Statistics, 1),
		ConfigChan:            make(chan<- Config, 1),
	}

	for _, aimConfig := range config.Aims {
		if !aimConfig.IsValid() {
			log.Println("Aim config not valid", aimConfig)
			continue
		}
		aim, err := NewAim(aimConfig)
		if err != nil {
			log.Println("Failed create aim", err)
			continue
		}
		pinger.Aims[aim.RemoteAddr.String()] = aim
		go aim.Start()
	}

	return &pinger, nil
}

func (p *Pinger) Start() {
	for {
		time.Sleep(p.CollectingStatsPeriod)

		statistics := Statistics{
			AimStats: make(map[string]AimStat),
		}

		for _, aim := range p.Aims {
			select {
			case stat := <-aim.StatChan:
				statistics.AimStats[aim.RemoteAddr.String()] = stat
			default:
				log.Println("Statistic for target", aim.RemoteAddr.String(), "not collected")
			}
		}

		statistics.Print()
	}
}
