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
	var wg sync.WaitGroup

    // coz we are netstalkers, not DoSers
    portRateLimiter := make(chan struct{}, 8)

	if len(s.Ports) == 0 {
		s.Ports = []int{0}
	}

	wg.Add(len(s.Ports))

	for _, port := range s.Ports {
		addr := net.TCPAddr{IP: ip, Port: port}
        portRateLimiter <- struct{}{}
        go func(){
            s.ServiceInterface.ScanAddr(addr, ch, &wg)
            <-portRateLimiter
        }()
	}

	wg.Wait()
	swg.Done() // w/o defer to speedup frequent calls
}
