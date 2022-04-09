package services

import (
	"net"
	"sync"
)

type ServiceInterface interface {
	ScanAddr(net.TCPAddr, chan<- HostResult, *sync.WaitGroup)
	Check(net.IP, chan<- HostResult, *sync.WaitGroup)
}

type Service struct {
	ServiceInterface
	Ports []int
}

var _ ServiceInterface = &Service{}

func (s *Service) Check(ip net.IP, ch chan<- HostResult, swg *sync.WaitGroup) {
	if len(s.Ports) == 0 {
		s.Ports = []int{0}
	}
	defer swg.Done()
	var wg sync.WaitGroup
	for _, port := range s.Ports {
		addr := net.TCPAddr{IP: ip, Port: port}
		wg.Add(1)
		go s.ServiceInterface.ScanAddr(addr, ch, &wg)
	}
	wg.Wait()
}
