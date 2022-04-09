package services

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fagci/gons/network"
	"github.com/fagci/gons/protocol"
	"github.com/fagci/gons/utils"
)

type HTTPService struct {
	Ports       []int
	connTimeout time.Duration
	paths       []string
	headerReg   *regexp.Regexp
	bodyReg     *regexp.Regexp
	client      *http.Client
}

var _ Service = &HTTPService{}

func NewHTTPService(ports []int, connTimeout time.Duration, paths []string, headerReg, bodyReg string) *HTTPService {
	if len(ports) == 0 {
		ports = []int{80, 443}
	}

	svc := HTTPService{
		Ports:       ports,
		connTimeout: connTimeout,
		paths:       paths,
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig.InsecureSkipVerify = true
	transport.DialContext = network.DialContextFunc(svc.connTimeout)
	transport.DisableKeepAlives = true // less memory leak, less errors

	svc.client = &http.Client{
		Transport: transport,
		Timeout:   protocol.RW_TIMEOUT,
	}

	if headerReg != "" {
		hReg, err := regexp.Compile(headerReg)
		if err != nil {
			panic("Bad headerReg: " + err.Error())
		}
		svc.headerReg = hReg
	}

	if bodyReg != "" {
		bReg, err := regexp.Compile(bodyReg)
		if err != nil {
			panic("Bad bodyReg: " + err.Error())
		}
		svc.bodyReg = bReg
	}

	return &svc
}

func (s *HTTPService) check(uri url.URL) (bool, [][]string, error) {
	response, err := s.client.Get(uri.String())
	if err != nil {
		return false, nil, err
	}

	defer response.Body.Close()

	if response.StatusCode > 400 {
		return false, nil, nil
	}

	for k, values := range response.Header {
		for _, v := range values {
			if s.headerReg != nil {
				if matches := s.headerReg.FindAllStringSubmatch(k+": "+v, -1); matches != nil {
					return true, matches, nil
				}
			}
		}
	}

	if response.ContentLength == -1 || response.ContentLength > 1024*1024 {
		return false, nil, nil
	}

	reader := io.LimitReader(response.Body, 1024*1024)
	body, err := io.ReadAll(reader)

	if err != nil {
		return false, nil, nil
	}

	if s.bodyReg != nil {
		if matches := s.bodyReg.FindAllStringSubmatch(string(body), -1); matches != nil {
			return true, matches, nil
		}
	}

	if s.headerReg == nil && s.bodyReg == nil {
		return true, nil, nil
	}

	return false, nil, nil
}

func (s *HTTPService) ScanAddr(addr net.TCPAddr, ch chan<- HostResult, wg *sync.WaitGroup) {
	defer wg.Done()
	scheme := "http"
	if addr.Port == 443 {
		scheme = "https"
	}
	for _, path := range s.paths {
		uri := url.URL{Scheme: scheme, Host: addr.String(), Path: path}
		ok, matches, err := s.check(uri)
		if err != nil {
			break
		}
		if ok {
			ch <- HostResult{
				Addr:    &addr,
				Details: &HTTPResult{Url: uri, Matches: matches},
			}
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
	Url     url.URL
	Matches [][]string
}

func (result *HTTPResult) ReplaceVars(cmd string) string {
	cmd = strings.ReplaceAll(cmd, "{result}", result.Url.String())
	cmd = strings.ReplaceAll(cmd, "{scheme}", result.Url.Scheme)
	cmd = strings.ReplaceAll(cmd, "{host}", result.Url.Host)
	cmd = strings.ReplaceAll(cmd, "{hostname}", result.Url.Hostname())
	cmd = strings.ReplaceAll(cmd, "{port}", result.Url.Port())
	cmd = strings.ReplaceAll(cmd, "{slug}", result.Slug())
	cmd = strings.ReplaceAll(cmd, "{matches_count}", fmt.Sprintf("%d", len(result.Matches)))
	// FIXME: unsafe
	/* for i, sm := range result.Matches {
		for j, m := range sm {
			escapedMatch := utils.FilterValueInQuotes(m)
			cmd = strings.ReplaceAll(cmd, fmt.Sprintf(`'{match_%d_%d}'`, i, j), escapedMatch)
			cmd = strings.ReplaceAll(cmd, fmt.Sprintf(`"{match_%d_%d}"`, i, j), escapedMatch)
		}
	} */
	return cmd
}

func (result *HTTPResult) Slug() string {
	return utils.Slugify(result.Url.String())
}

func (result *HTTPResult) String() string {
	if result.Matches == nil {
		return result.Url.String()
	}

	sb := strings.Builder{}
	for i, sm := range result.Matches {
		if len(sm) == 1 {
			continue
		}
		for j, m := range sm {
			if j > 0 {
				sb.WriteString(fmt.Sprintf("Match[%d][%d]: %s\n", i, j, m))
			}
		}
	}

	return result.Url.String() + "\n" + sb.String()
}
