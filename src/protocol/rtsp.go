package protocol

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fagci/gons/src/generators"
	"github.com/fagci/gons/src/network"
)

type RTSP struct {
	timeout time.Duration
	conn    net.Conn
	cseq    int
	paths   []string
	addr    string
}

const RW_TIMEOUT = time.Second * 2

const _RTSP_TPL = "%s %s RTSP/1.0\r\n" +
	"Accept: application/sdp\r\n" +
	"CSeq: %d\r\n" +
	"User-Agent: Lavf59.16.100\r\n\r\n"

func (r *RTSP) Request(req string) (int, error) {
	if e := r.conn.SetDeadline(time.Now().Add(RW_TIMEOUT)); e != nil {
		return 0, e
	}

	if _, e := r.conn.Write([]byte(req)); e != nil {
		return 0, e
	}

	responseBytes := make([]byte, 512)
	n, e := r.conn.Read(responseBytes)
	if e != nil {
		return 0, e
	}

	f := strings.Fields(string(responseBytes[:n]))
	if len(f) >= 2 && strings.HasPrefix(f[0], "RTSP") {
		return strconv.Atoi(f[1])
	}

	return 0, errors.New("Bad response")
}

func (r *RTSP) Query(path string) string {
	method := "DESCRIBE"

	if path == "*" {
		method = "OPTIONS"
	} else {
		path = fmt.Sprintf("rtsp://%s%s", r.addr, path)
	}

	r.cseq++
	return fmt.Sprintf(_RTSP_TPL, method, path, r.cseq)
}

func (r *RTSP) Check() (url.URL, error) {
	var err error
	var url url.URL

	if r.conn, err = network.DialTimeout("tcp", r.addr, r.timeout); err != nil {
		return url, err
	}

	defer r.conn.Close()

	if _, err = r.Request(r.Query("*")); err != nil {
		return url, err
	}

	var code int
	code, err = r.Request(r.Query(generators.RandomPath(6, 12)))
	if err != nil {
		return url, err
	}

	if code == 200 {
		code, err = r.Request("/")
		if err == nil && code == 200 {
			return r.result("/"), nil
		}
		return url, errors.New("RTSP: fake")
	}

	for _, path := range r.paths {
		code, err = r.Request(r.Query(path))
		if err != nil {
			return url, err
		}
		if code == 401 {
			return url, errors.New("RTSP: unauthorized")
		}
		if code == 200 {
			return r.result(path), nil
		}
	}

	return url, errors.New("RTSP: bad response")
}

func (r *RTSP) result(path string) url.URL {
	return url.URL{
		Scheme: "rtsp",
		Host:   r.addr,
		Path:   path,
	}
}

func NewRTSP(addr net.Addr, paths []string, timeout time.Duration) *RTSP {
	return &RTSP{
		timeout: timeout,
		paths:   paths,
		addr:    addr.String(),
	}
}
