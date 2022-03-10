package pinger

import (
	"log"
	"math/rand"
	"net"
	"time"
	"udp-ping/packet"
)

type Aim struct {
	Name          string
	PingPeriod    time.Duration
	Conn          *net.UDPConn
	LocalAddr     net.UDPAddr
	RemoteAddr    net.UDPAddr
	Stat          AimStat
	StatChan      chan AimStat
	SendPacket    chan packet.Packet
	RecvPacket    chan packet.Packet
	PacketIndexes map[Index]time.Time
	StopChan      chan struct{}
	ReconnectChan chan struct{}
}

type Index struct {
	int
}

type AimStat struct {
	Name             string
	CreationTime     time.Time
	LocalAddr        net.UDPAddr
	RemoteAddr       net.UDPAddr
	SendPacketsCount int
	RecvPacketsCount int
	RttSum           time.Duration
	AverageRtt       time.Duration
}

func NewAim(aimConfig AimConfig) (*Aim, error) {
	aim := &Aim{
		Name:          aimConfig.Name,
		PacketIndexes: make(map[Index]time.Time),
		StatChan:      make(chan AimStat, 1),
		SendPacket:    make(chan packet.Packet, 1),
		RecvPacket:    make(chan packet.Packet, 1),
		StopChan:      make(chan struct{}, 1),
		ReconnectChan: make(chan struct{}, 1),

		Stat: AimStat{
			Name:             aimConfig.Name,
			CreationTime:     time.Now(),
			SendPacketsCount: 0,
			RecvPacketsCount: 0,
			RttSum:           0,
			AverageRtt:       0,
		},
	}

	if localAddr, err := net.ResolveUDPAddr("udp", aimConfig.LocalAddr); err != nil {
		return nil, err
	} else {
		aim.LocalAddr = *localAddr
		aim.Stat.LocalAddr = *localAddr
	}

	if remoteAddr, err := net.ResolveUDPAddr("udp", aimConfig.RemoteAddr); err != nil {
		return nil, err
	} else {
		aim.RemoteAddr = *remoteAddr
		aim.Stat.RemoteAddr = *remoteAddr
	}

	if err := aim.connect(); err != nil {
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
					log.Println("packet index not found", pack.Index)
				} else {
					delete(aim.PacketIndexes, Index{pack.Index})
					aim.Stat.RecvPacketsCount++
					packRtt := time.Since(sendTime)
					aim.Stat.RttSum += packRtt
					aim.Stat.AverageRtt = aim.Stat.RttSum / time.Duration(aim.Stat.RecvPacketsCount)
					// log.Println("OK", packRtt, aim.Stats.AverageRtt)
					select {
					case aim.StatChan <- aim.Stat:
					default:
					}
				}

			case <-aim.StopChan:
				log.Println("Stoping ping aim:", aim.RemoteAddr.String())
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
				log.Println(err)
				return
			}
			incomePacket := packet.Unmarshal(packet.ToPacketBin(buffer[:n]))
			if incomePacket.Type == packet.ECHO_REPLY {
				aim.RecvPacket <- incomePacket
			} else {
				log.Println("not echo packet")
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
			log.Println(err)
			return
		}

		time.Sleep(aim.PingPeriod)
	}
}

func (aim *Aim) connect() error {
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
