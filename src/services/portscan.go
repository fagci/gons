package services

import (
	"gons/src/models"
	"net"
	"sync"
	"time"
)

type PortscanService struct {
	Ports   []int
	timeout time.Duration
}

func NewPortscanService(ports []int, timeout time.Duration) *PortscanService {
	return &PortscanService{Ports: ports, timeout: timeout}
}

func (s *PortscanService) ScanAddr(addr net.TCPAddr, ch *chan models.HostResult, wg *sync.WaitGroup) {
    defer wg.Done()
	if conn, err := net.DialTimeout("tcp", addr.String(), s.timeout); err == nil {
		conn.Close()
		*ch <- models.HostResult{Addr: &addr}
	}
}

func (s *PortscanService) Check(ip net.IP) <-chan models.HostResult {
    var wg sync.WaitGroup
	ch := make(chan models.HostResult)
	go func() {
		defer close(ch)
		for _, port := range s.Ports {
			addr := net.TCPAddr{IP: ip, Port: port}
            wg.Add(1)
            go s.ScanAddr(addr, &ch, &wg)
		}
        wg.Wait()
	}()
	return ch
}
