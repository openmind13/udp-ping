package aim

import (
	"crypto/rand"
	"net"
	"time"
	"udp-pinger/packet"

	"github.com/sirupsen/logrus"
)

type Aim struct {
	Name          string
	Config        AimConfig
	PingPeriod    time.Duration
	Conn          *net.UDPConn
	LocalAddr     net.UDPAddr
	RemoteAddr    net.UDPAddr
	Stat          Stat
	StatChan      chan Stat
	SendPacket    chan packet.Packet
	RecvPacket    chan packet.Packet
	PacketIndexes map[Index]time.Time
	StopChan      chan struct{}
	ReconnectChan chan struct{}
	ConfigChan    chan AimConfig
}

type Index struct {
	int
}

type Stat struct {
	Name             string
	CreationTime     time.Time
	LocalAddr        net.UDPAddr
	RemoteAddr       net.UDPAddr
	SendPacketsCount int
	RecvPacketsCount int
	RttSum           time.Duration
	AverageRtt       time.Duration
}

func New(aimConfig AimConfig, pingPeriod time.Duration) (*Aim, error) {
	aim := &Aim{
		Name:          aimConfig.Name,
		Config:        aimConfig,
		PingPeriod:    pingPeriod,
		PacketIndexes: make(map[Index]time.Time),
		StatChan:      make(chan Stat, 1),
		SendPacket:    make(chan packet.Packet, 1),
		RecvPacket:    make(chan packet.Packet, 1),
		StopChan:      make(chan struct{}, 1),
		ReconnectChan: make(chan struct{}, 1),
		ConfigChan:    make(chan AimConfig, 1),

		Stat: Stat{
			Name:             aimConfig.Name,
			CreationTime:     time.Now(),
			SendPacketsCount: 0,
			RecvPacketsCount: 0,
			RttSum:           0,
			AverageRtt:       0,
		},
	}

	if err := aim.reconnect(); err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case pack := <-aim.SendPacket:
				aim.PacketIndexes[Index{pack.Index}] = pack.Time
				aim.Stat.SendPacketsCount++

			case pack := <-aim.RecvPacket:
				sendTime, ok := aim.PacketIndexes[Index{pack.Index}]
				if !ok {
					logrus.Warn("packet index not found", pack.Index)
				} else {
					delete(aim.PacketIndexes, Index{pack.Index})
					aim.Stat.RecvPacketsCount++
					packRtt := time.Since(sendTime)
					aim.Stat.RttSum += packRtt
					aim.Stat.AverageRtt = aim.Stat.RttSum / time.Duration(aim.Stat.RecvPacketsCount)
					select {
					case aim.StatChan <- aim.Stat:
					default:
					}
				}

			case newConfig := <-aim.ConfigChan:
				logrus.Debug("from aim new config")
				if aim.Config == newConfig {
					break
				}
				if aim.Config.LocalAddr != newConfig.LocalAddr || aim.Config.RemoteAddr != newConfig.RemoteAddr {
					// reconnect not needed
					aim.Conn.Close()
					if err := aim.reconnect(); err != nil {
						logrus.WithFields(logrus.Fields{
							"aim_name": aim.Name,
							"error":    err.Error(),
						}).Error("Failed to reconnect")
						break
					}
				} else {
					logrus.Debug("Addresses doesn't changed. ", aim.Name)
				}

			case <-aim.StopChan:
				logrus.Info("Stoping ping aim:", aim.RemoteAddr.String())
				aim.Conn.Close()
				return
			}
		}
	}()

	return aim, nil
}

func (aim *Aim) Start() {
	go func() {
		for {
			buffer := make([]byte, packet.MAX_SAFE_UDP_PAYLOAD_SIZE_BYTES)
			n, _, err := aim.Conn.ReadFromUDP(buffer)
			if err != nil {
				logrus.Error(err)
				return
			}
			incomePacket := packet.Unmarshal(packet.ToPacketBin(buffer[:n]))
			if incomePacket.Type == packet.ECHO_REPLY {
				aim.RecvPacket <- incomePacket
			} else {
				logrus.Warn("not echo packet")
			}
		}
	}()

	for {
		index := rand.Int()
		echoPacket := packet.Packet{
			Index: index,
			Type:  packet.ECHO_REQUEST,
			Time:  time.Now(),
		}
		aim.SendPacket <- echoPacket
		_, err := aim.Conn.WriteToUDP(echoPacket.Marshal().ToSlice(), &aim.RemoteAddr)
		if err != nil {
			logrus.Error(err)
			return
		}

		time.Sleep(aim.PingPeriod)
	}
}

func (aim *Aim) reconnect() error {
	if localAddr, err := net.ResolveUDPAddr("udp", aim.Config.LocalAddr); err != nil {
		return err
	} else {
		aim.LocalAddr = *localAddr
		aim.Stat.LocalAddr = *localAddr
	}

	if remoteAddr, err := net.ResolveUDPAddr("udp", aim.Config.RemoteAddr); err != nil {
		return err
	} else {
		aim.RemoteAddr = *remoteAddr
		aim.Stat.RemoteAddr = *remoteAddr
	}

	conn, err := net.ListenUDP("udp", &aim.LocalAddr)
	if err != nil {
		return err
	}
	aim.Conn = conn
	return nil
}

func (aim *Aim) Stop() {
	aim.StopChan <- struct{}{}
}

func (aim *Aim) ChangeConfig() error {

	return nil
}
