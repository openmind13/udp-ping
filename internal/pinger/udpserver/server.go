package udpserver

import "net"

type UdpServer struct {
	socket *net.UDPConn
}

func New(listenAddr string) (*UdpServer, error) {
	s := &UdpServer{}
	resolvedAddr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenUDP("udp", resolvedAddr)
	if err != nil {
		return nil, err
	}
	s.socket = listener
	return s, nil
}
