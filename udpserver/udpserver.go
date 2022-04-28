package main

import (
	"flag"
	"net"
	"time"
	"udp-pinger/packet"

	"github.com/sirupsen/logrus"
)

var (
	localAddrFlag = flag.String("addr", "0.0.0.0:9000", "local udp addr")
)

const (
	STATS_PERIOD              = 2 * time.Second
	CHECK_CLIENTS_PERIOD      = 10 * time.Second
	CLIENT_DISCONNECT_TIMEOUT = 15 * time.Second
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "15:04:05 02-01-2006",
		DisableColors:   false,
		FullTimestamp:   true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	flag.Parse()

	localAddr, err := net.ResolveUDPAddr("udp", *localAddrFlag)
	if err != nil {
		logrus.Fatal(err)
	}

	udpServer, err := NewServer(localAddr)
	if err != nil {
		logrus.Fatal(err)
	}
	go udpServer.Start()

	for {
		statistic := <-udpServer.StatisticChan
		statistic.Print()
	}
}

type Client struct {
	Addr               net.UDPAddr
	ConnectionTime     time.Time
	LastRecvPacketTime time.Time
	Stat               ClientStat
}

type UDPServer struct {
	conn               *net.UDPConn
	clients            map[string]Client
	recvPacketInfoChan chan PacketInfo
	sendPacketInfoChan chan PacketInfo
	// checkClientsHook   chan struct{}
	// statisticHookChan  chan struct{}
	StatisticChan chan Statistic
}

type ClientStat struct {
	EchoRequestPacketCount int
	EchoReplyPacketCount   int
	// AnotherPacketCount     int
}

type Statistic struct {
	Time time.Time
	Data map[string]ClientStat
}

func (s Statistic) Print() {
	logrus.Println("Statistic on", s.Time.Round(time.Second).String())
	for addr, stat := range s.Data {
		logrus.Println("Addr:", addr, " ping:", stat.EchoRequestPacketCount, " pong:", stat.EchoReplyPacketCount)
	}
	logrus.Println("-------------------------------------------")
}

type PacketInfo struct {
	Packet packet.Packet
	Addr   net.UDPAddr
}

func NewServer(addr *net.UDPAddr) (*UDPServer, error) {
	s := &UDPServer{
		clients:            make(map[string]Client),
		recvPacketInfoChan: make(chan PacketInfo, 1),
		sendPacketInfoChan: make(chan PacketInfo, 1),
		// checkClientsHook:   make(chan struct{}, 1),
		// statisticHookChan: make(chan struct{}, 1),
		StatisticChan: make(chan Statistic, 1),
	}

	logrus.Println("Starting listening udp on:", addr)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	s.conn = conn

	// go func() {
	// 	for {
	// 		time.Sleep(STATS_PERIOD)
	// 		s.statisticHookChan <- struct{}{}
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		time.Sleep(CHECK_CLIENTS_PERIOD)
	// 		s.checkClientsHook <- struct{}{}
	// 	}
	// }()

	statisticHook := time.NewTicker(STATS_PERIOD)
	checkClientHook := time.NewTicker(CHECK_CLIENTS_PERIOD)

	go func() {
		for {
			select {
			case <-statisticHook.C:
			case packInfo := <-s.recvPacketInfoChan:
				client, ok := s.clients[packInfo.Addr.String()]
				if !ok {
					s.clients[packInfo.Addr.String()] = Client{
						Addr:           packInfo.Addr,
						ConnectionTime: time.Now(),
						Stat: ClientStat{
							EchoRequestPacketCount: 1,
							EchoReplyPacketCount:   0,
						},
					}
				} else {
					client.Stat.EchoRequestPacketCount++
					client.LastRecvPacketTime = time.Now()
					s.clients[client.Addr.String()] = client
				}
			case packInfo := <-s.sendPacketInfoChan:
				client, ok := s.clients[packInfo.Addr.String()]
				if !ok {
					s.clients[packInfo.Addr.String()] = Client{
						Addr: packInfo.Addr,
					}
				} else {
					client.Stat.EchoReplyPacketCount++
					s.clients[packInfo.Addr.String()] = client
				}
			case <-statisticHook.C:
				statistic := Statistic{
					Data: make(map[string]ClientStat),
				}
				for addr, client := range s.clients {
					statistic.Data[addr] = client.Stat
				}
				statistic.Time = time.Now()
				select {
				case s.StatisticChan <- statistic:
				default: // ignore if we chan is full
				}
			case <-checkClientHook.C:
				disconnectKeyCandidatesList := []string{}
				for addrStr, client := range s.clients {
					if time.Since(client.LastRecvPacketTime) > CLIENT_DISCONNECT_TIMEOUT {
						disconnectKeyCandidatesList = append(disconnectKeyCandidatesList, addrStr)
					}
				}
				for _, addr := range disconnectKeyCandidatesList {
					delete(s.clients, addr)
				}
			}
		}
	}()

	return s, nil
}

func (s *UDPServer) Start() {
	for {
		buffer := make([]byte, packet.MAX_SAFE_UDP_PAYLOAD_SIZE_BYTES)

		_, remoteAddr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			logrus.Fatal(err)
		}

		incomingPacket := packet.Unmarshal(packet.ToPacketBin(buffer))
		if incomingPacket.Type != packet.ECHO_REQUEST {
			logrus.Warn("not echo request")
			continue
		}
		s.recvPacketInfoChan <- PacketInfo{
			Packet: incomingPacket,
			Addr:   *remoteAddr,
		}

		echoReplyPacket := packet.Packet{
			Index: incomingPacket.Index,
			Time:  incomingPacket.Time,
			Type:  packet.ECHO_REPLY,
		}
		_, err = s.conn.WriteToUDP(echoReplyPacket.Marshal().ToSlice(), remoteAddr)
		if err != nil {
			logrus.Fatal(err)
		}
		s.sendPacketInfoChan <- PacketInfo{
			Packet: echoReplyPacket,
			Addr:   *remoteAddr,
		}
	}
}
