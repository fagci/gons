package services

import (
	"net"
	"sync"
)

type DummyService struct {
}

func NewDummyService() *DummyService {
	return &DummyService{}
}

func (ds *DummyService) Check(ip net.IP, ch chan<- HostResult, wg *sync.WaitGroup) {
	defer wg.Done()
	ch <- HostResult{Addr: &net.TCPAddr{IP: ip}}
}
