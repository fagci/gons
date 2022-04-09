package services

import (
	"net"
	"sync"
)

type DummyService struct {
	*Service
}

func NewDummyService() *DummyService {
    s := &DummyService{Service: &Service{}}
	s.ServiceInterface = interface{}(s).(ServiceInterface)
	return s
}

func (s *Service) ScanAddr(addr net.TCPAddr, ch chan<- HostResult, wg *sync.WaitGroup) {
	defer wg.Done()
	ch <- HostResult{Addr: &addr}
}
