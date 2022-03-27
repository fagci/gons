package services

import (
	"gons/src/models"
	"net"
)

type DummyService struct {
}

func NewDummyService() *DummyService {
	return &DummyService{}
}

func (ds *DummyService) Check(ip net.IP) <-chan models.HostResult {
	ch := make(chan models.HostResult)
	go func() {
        defer close(ch)
		ch <- models.HostResult{
			Addr: &net.TCPAddr{IP: ip},
		}
	}()
	return ch
}
