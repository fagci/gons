package services

import (
	"gons/src/models"
	"net"
	"time"
)

type PortscanService struct {
	Ports   []int
	timeout time.Duration
}

func NewPortscanService(ports []int, timeout time.Duration) *PortscanService {
	return &PortscanService{Ports: ports, timeout: timeout}
}

func (s *PortscanService) Check(ip net.IP) <-chan models.HostResult {
	ch := make(chan models.HostResult)
	go func() {
		defer close(ch)
		for _, port := range s.Ports {
			addr := net.TCPAddr{IP: ip, Port: port}
			if conn, err := net.DialTimeout("tcp", addr.String(), s.timeout); err == nil {
				conn.Close()
				ch <- models.HostResult{
					Addr: &addr,
				}
			}
		}
	}()
	return ch
}
