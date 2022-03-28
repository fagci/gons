package services

import (
	"gons/src/models"
	"gons/src/protocol"
	"net"
	"sync"
	"time"
)

type RTSPService struct {
	Ports       []int
	connTimeout time.Duration
	paths       []string
}

func NewRTSPService(ports []int, connTimeout time.Duration, paths []string) *RTSPService {
	return &RTSPService{
		Ports:       ports,
		connTimeout: connTimeout,
		paths:       paths,
	}
}
func (s *RTSPService) ScanAddr(addr net.TCPAddr, ch *chan models.HostResult, wg *sync.WaitGroup) {
	defer wg.Done()
	r := protocol.NewRTSP(&addr, s.paths, s.connTimeout)
	for res := range r.Check() {
		*ch <- res
	}
}

func (s *RTSPService) Check(host net.IP) <-chan models.HostResult {
	var wg sync.WaitGroup
	ch := make(chan models.HostResult)
	go func() {
		defer close(ch)
		for _, port := range s.Ports {
			addr := net.TCPAddr{IP: host, Port: port}
			wg.Add(1)
			go s.ScanAddr(addr, &ch, &wg)
		}
		wg.Wait()
	}()
	return ch
}
