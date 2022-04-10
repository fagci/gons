package services

import (
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/fagci/gonr/generators"
	"github.com/fagci/gons/protocol"
	"github.com/fagci/gons/utils"
)

type RTSPService struct {
	*Service
	Ports    []int
	timeout  time.Duration
	paths    []string
	fakePath string
}

func NewRTSPService(ports []int, timeout time.Duration, paths []string) *RTSPService {
	if len(ports) == 0 {
		ports = append(ports, 554)
	}
	s := &RTSPService{
		timeout:  timeout,
		paths:    paths,
		Service:  &Service{Ports: ports},
		fakePath: generators.RandomPath(6, 12),
	}
	s.ServiceInterface = interface{}(s).(ServiceInterface)
	return s
}
func (s *RTSPService) ScanAddr(addr net.TCPAddr, ch chan<- HostResult, wg *sync.WaitGroup) {
	defer wg.Done()
	r := protocol.NewRTSP(&addr, s.paths, s.fakePath, s.timeout)
	if res, err := r.Check(); err == nil {
		ch <- HostResult{
			Addr:    &addr,
			Details: &RTSPResult{Url: res},
		}
	}
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
