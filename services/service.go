package services

import (
	"net"
	"sync"
)

type Service interface {
	Check(net.IP, chan<- HostResult, *sync.WaitGroup)
}
