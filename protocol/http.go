package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fagci/gonr/generators"
	"github.com/fagci/gons/network"
)

type HTTP struct {
	timeout            time.Duration
	conn               net.Conn
	paths              []string
	addr               string
	headerReg, bodyReg *regexp.Regexp
}

const _HTTP_TPL = "%s %s HTTP/1.0\r\n" +
	"User-Agent: Mozilla/5.0\r\n\r\n"

func (r *HTTP) Request(req string) (int, error) {
	if e := r.conn.SetDeadline(time.Now().Add(RW_TIMEOUT)); e != nil {
		return 0, e
	}

	if _, e := r.conn.Write([]byte(req)); e != nil {
		return 0, e
	}

	reader := bufio.NewReader(r.conn)
	tp := textproto.NewReader(reader)

	firstLine, err := tp.ReadLine()
	if err != nil {
		return 0, err
	}

	f := strings.Fields(firstLine)
	if len(f) < 2 || !strings.HasPrefix(f[0], "HTTP") {
		return 0, errors.New("Bad response")
	}
	code, err := strconv.Atoi(f[1])
	if err != nil || code != 200 {
		return code, err
	}

	isHeader := true
	for {
		line, err := tp.ReadLine()
		if err != nil {
			break
		}
		if line == "" {
			isHeader = false
		}
		if isHeader && r.headerReg != nil {
			if r.headerReg.MatchString(line) {
				return 200, nil
			}
		}
		if !isHeader && r.bodyReg != nil {
			if r.bodyReg.MatchString(line) {
				return 200, nil
			}
		}
	}

	return code, nil
}

func (r *HTTP) Query(path string) string {
	method := "GET"

	path = fmt.Sprintf("http://%s%s", r.addr, path)

	return fmt.Sprintf(_HTTP_TPL, method, path)
}

func (r *HTTP) Check() (url.URL, error) {
	var err error
	var url url.URL

	if r.conn, err = network.DialTimeout("tcp", r.addr, r.timeout); err != nil {
		return url, err
	}

	defer r.conn.Close()

	var code int
	code, err = r.Request(r.Query(generators.RandomPath(6, 12)))
	if err != nil {
		return url, err
	}

	if code == 200 {
		return url, errors.New("HTTP: fake")
	}

	for _, path := range r.paths {
		code, err = r.Request(r.Query(path))
		if err != nil {
			return url, err
		}
		if code == 401 {
			return url, errors.New("HTTP: unauthorized")
		}
		if code == 200 {
			return r.result(path), nil
		}
	}

	return url, errors.New("HTTP: bad response")
}

func (r *HTTP) result(path string) url.URL {
	return url.URL{
		Scheme: "http",
		Host:   r.addr,
		Path:   path,
	}
}

func NewHTTP(addr net.Addr, paths []string, timeout time.Duration, headerReg, bodyReg *regexp.Regexp) *HTTP {
	return &HTTP{
		timeout:   timeout,
		paths:     paths,
		addr:      addr.String(),
		headerReg: headerReg,
		bodyReg:   bodyReg,
	}
}
