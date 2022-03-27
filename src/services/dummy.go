package services

import (
	"go-ns/src/models"
	"net"
)

type DummyService struct {
	ch chan models.HostResult
}

func NewDummyService() *DummyService {
	return &DummyService{
        ch: make(chan models.HostResult,1),
    }
}

func (ds *DummyService) Check(ip net.IP) <-chan models.HostResult {
	ds.ch <- models.HostResult{
        Addr: &net.TCPAddr{
            IP: ip,
        },
    }
	return ds.ch
}
