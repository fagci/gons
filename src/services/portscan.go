package services

import (
	"github.com/fagci/gons/src/network"
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

func (s *PortscanService) ScanAddr(addr net.TCPAddr, ch chan<- HostResult, wg *sync.WaitGroup) {
	defer wg.Done()
	if conn, err := network.DialTimeout("tcp", addr.String(), s.timeout); err == nil {
		conn.Close()
		ch <- HostResult{Addr: &addr}
	}
}

func (s *PortscanService) Check(ip net.IP, ch chan<- HostResult, swg *sync.WaitGroup) {
    defer swg.Done()
	var wg sync.WaitGroup
	go func() {
		for _, port := range s.Ports {
			addr := net.TCPAddr{IP: ip, Port: port}
			wg.Add(1)
			go s.ScanAddr(addr, ch, &wg)
		}
		wg.Wait()
	}()
}
