package services

import (
	"go-ns/src/models"
	"go-ns/src/protocol"
	"time"
)

type RTSPService struct {
	Port  int
	paths []string
}

func NewRTSPService(port int, paths []string) *RTSPService {
	return &RTSPService{
		Port:  port,
		paths: paths,
	}
}

func (rs *RTSPService) Check(host string) <-chan models.Result {
	r := protocol.NewRTSP(host, rs.Port, rs.paths, time.Second*2)
	return r.Check()
}

/* func NewRTSPService2(port int, paths []string) func (string) <-chan models.Result {
    rs := &RTSPService{
		Port:  port,
		paths: paths,
	}
	return func(host string) <-chan models.Result {
		r := protocol.NewRTSP(host, rs.Port, rs.paths, time.Second*2)
		return r.Check()
	}
} */
