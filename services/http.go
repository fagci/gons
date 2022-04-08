package services

import (
	"net"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fagci/gons/protocol"
	"github.com/fagci/gons/utils"
)

type HTTPService struct {
	Ports       []int
	connTimeout time.Duration
	paths       []string
	headerReg   *regexp.Regexp
	bodyReg     *regexp.Regexp
}

var _ Service = &HTTPService{}

func NewHTTPService(ports []int, connTimeout time.Duration, paths []string, headerReg, bodyReg string) *HTTPService {
	if len(ports) == 0 {
		ports = append(ports, 554)
	}
    svc := HTTPService{
		Ports:       ports,
		connTimeout: connTimeout,
		paths:       paths,
	}

    if headerReg != "" {
        hReg, err := regexp.Compile(headerReg)
        if err !=nil {
            panic("Bad headerReg: " + err.Error())
        }
        svc.headerReg = hReg
    }

    if bodyReg != "" {
        bReg, err := regexp.Compile(bodyReg)
        if err !=nil {
            panic("Bad bodyReg: " + err.Error())
        }
        svc.bodyReg = bReg
    }

    return &svc
}
func (s *HTTPService) ScanAddr(addr net.TCPAddr, ch chan<- HostResult, wg *sync.WaitGroup) {
	defer wg.Done()
	r := protocol.NewHTTP(&addr, s.paths, s.connTimeout, s.headerReg, s.bodyReg)
	if res, err := r.Check(); err == nil {
		ch <- HostResult{
			Addr:    &addr,
			Details: &HTTPResult{Url: res},
		}
	}
}

func (s *HTTPService) Check(host net.IP, ch chan<- HostResult, swg *sync.WaitGroup) {
	defer swg.Done()
	var wg sync.WaitGroup
	for _, port := range s.Ports {
		addr := net.TCPAddr{IP: host, Port: port}
		wg.Add(1)
		go s.ScanAddr(addr, ch, &wg)
	}
	wg.Wait()
}

type HTTPResult struct {
	Url url.URL
}

func (result *HTTPResult) ReplaceVars(cmd string) string {
	cmd = strings.ReplaceAll(cmd, "{result}", result.Url.String())
	cmd = strings.ReplaceAll(cmd, "{scheme}", result.Url.Scheme)
	cmd = strings.ReplaceAll(cmd, "{host}", result.Url.Host)
	cmd = strings.ReplaceAll(cmd, "{hostname}", result.Url.Hostname())
	cmd = strings.ReplaceAll(cmd, "{port}", result.Url.Port())
	cmd = strings.ReplaceAll(cmd, "{slug}", result.Slug())
	return cmd
}

func (result *HTTPResult) Slug() string {
	return utils.Slugify(result.Url.String())
}

func (result *HTTPResult) String() string {
	return result.Url.String()
}
