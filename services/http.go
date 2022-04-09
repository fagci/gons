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
	*Service
	timeout   time.Duration
	paths     []string
	headerReg *regexp.Regexp
	bodyReg   *regexp.Regexp
	client    *http.Client
}

const MAX_HTTP_BODY_LENGTH = 1024 * 1024

func NewHTTPService(ports []int, timeout time.Duration, paths []string, headerReg, bodyReg string) *HTTPService {
	if len(ports) == 0 {
		ports = []int{80, 443}
	}

	s := &HTTPService{
		timeout: timeout,
		paths:   paths,
		Service: &Service{Ports: ports},
	}
	s.ServiceInterface = interface{}(s).(ServiceInterface)

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig.InsecureSkipVerify = true
	transport.DialContext = network.DialContextFunc(s.timeout)
	transport.DisableKeepAlives = true // less memory leak, less errors

	s.client = &http.Client{
		Transport: transport,
		Timeout:   protocol.RW_TIMEOUT,
	}

	if headerReg != "" {
		hReg, err := regexp.Compile(headerReg)
		if err != nil {
			panic("Bad headerReg: " + err.Error())
		}
		s.headerReg = hReg
	}

	if bodyReg != "" {
		bReg, err := regexp.Compile(bodyReg)
		if err != nil {
			panic("Bad bodyReg: " + err.Error())
		}
		s.bodyReg = bReg
	}

	return s
}

func (s *HTTPService) check(uri url.URL) (bool, [][]string, error) {
	var matches [][]string
	response, err := s.client.Get(uri.String())
	if err != nil {
		return false, nil, err
	}

	defer response.Body.Close()

	if response.StatusCode > 400 {
		return false, nil, nil
	}

	if s.headerReg != nil {
		for k, values := range response.Header {
			for _, v := range values {
				matches = append(matches, s.headerReg.FindAllStringSubmatch(k+": "+v, -1)...)
			}
		}
	}

	if s.bodyReg != nil {
		body, err := io.ReadAll(io.LimitReader(response.Body, MAX_HTTP_BODY_LENGTH))
		if err != nil {
			return false, matches, nil
		}

		matches = append(matches, s.bodyReg.FindAllStringSubmatch(string(body), -1)...)
	}

	return (s.headerReg == nil && s.bodyReg == nil) || len(matches) != 0, matches, nil
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

	if sb.Len() == 0 {
		return result.Url.String()
	}

	return result.Url.String() + "\n" + sb.String()
}
