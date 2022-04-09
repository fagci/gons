package services

import (
	"net"
	"sync"
	"time"

	"github.com/fagci/gons/network"
)

type PortscanService struct {
	*Service
	timeout time.Duration
}

func NewPortscanService(ports []int, timeout time.Duration) *PortscanService {
	s := &PortscanService{
		timeout: timeout,
		Service: &Service{Ports: ports},
	}
	s.ServiceInterface = interface{}(s).(ServiceInterface)
	return s
}

func (s *PortscanService) ScanAddr(addr net.TCPAddr, ch chan<- HostResult, wg *sync.WaitGroup) {
	defer wg.Done()
	if conn, err := network.DialTimeout("tcp", addr.String(), s.timeout); err == nil {
		conn.Close()
		ch <- HostResult{Addr: &addr}
	}
}
