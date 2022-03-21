package services

import "go-ns/src/protocol"

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

func (rs *RTSPService) Check(host string) <-chan string {
	r := protocol.NewRTSP(host, rs.Port, rs.paths)
	return r.Check()
}
