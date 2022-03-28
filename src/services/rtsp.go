package services

import (
	"gons/src/models"
	"gons/src/protocol"
	"net"
	"time"
)

type RTSPService struct {
	Port        int
	paths       []string
	connTimeout time.Duration
}

func NewRTSPService(port int, paths []string, connTimeout time.Duration) *RTSPService {
	return &RTSPService{
		Port:        port,
		paths:       paths,
		connTimeout: connTimeout,
	}
}

func (rs *RTSPService) Check(host net.IP) <-chan models.HostResult {
	addr := net.TCPAddr{IP: host, Port: rs.Port}
	r := protocol.NewRTSP(&addr, rs.paths, rs.connTimeout)
	return r.Check()
}
