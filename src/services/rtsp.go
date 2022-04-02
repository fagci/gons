package services

import (
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/fagci/gons/src/protocol"
	"github.com/fagci/gons/src/utils"
)

type RTSPService struct {
	Ports       []int
	connTimeout time.Duration
	paths       []string
}

var _ Service = &RTSPService{}

func NewRTSPService(ports []int, connTimeout time.Duration, paths []string) *RTSPService {
	if len(ports) == 0 {
		ports = append(ports, 554)
	}
	return &RTSPService{
		Ports:       ports,
		connTimeout: connTimeout,
		paths:       paths,
	}
}
func (s *RTSPService) ScanAddr(addr net.TCPAddr, ch chan<- HostResult, wg *sync.WaitGroup) {
	defer wg.Done()
	r := protocol.NewRTSP(&addr, s.paths, s.connTimeout)
	if res, err := r.Check(); err == nil {
		ch <- HostResult{
			Addr:    &addr,
			Details: &RTSPResult{Url: res},
		}
	}
}

func (s *RTSPService) Check(host net.IP, ch chan<- HostResult, swg *sync.WaitGroup) {
	defer swg.Done()
	var wg sync.WaitGroup
	for _, port := range s.Ports {
		addr := net.TCPAddr{IP: host, Port: port}
		wg.Add(1)
		go s.ScanAddr(addr, ch, &wg)
	}
	wg.Wait()
}

type RTSPResult struct {
	Url url.URL
}

func (result *RTSPResult) ReplaceVars(cmd string) string {
	cmd = strings.ReplaceAll(cmd, "{result}", result.Url.String())
	cmd = strings.ReplaceAll(cmd, "{scheme}", result.Url.Scheme)
	cmd = strings.ReplaceAll(cmd, "{host}", result.Url.Host)
	cmd = strings.ReplaceAll(cmd, "{hostname}", result.Url.Hostname())
	cmd = strings.ReplaceAll(cmd, "{port}", result.Url.Port())
	cmd = strings.ReplaceAll(cmd, "{slug}", result.Slug())
	return cmd
}

func (result *RTSPResult) Slug() string {
	return utils.Slugify(result.Url.String())
}

func (result *RTSPResult) String() string {
	return result.Url.String()
}
