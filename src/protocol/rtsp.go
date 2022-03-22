package protocol

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type RTSP struct {
	Host    string
	Port    int
	timeout time.Duration
	conn    net.Conn
	cseq    int
	paths   []string
	ch      chan string
}

const _RTSP_TPL = "%s %s RTSP/1.0\r\n" +
	"CSeq: %d\r\n" +
	"User-Agent: LibVLC/3.0.0\r\n" +
	"Accept: application/sdp\r\n\r\n"

func (r *RTSP) Request(req string) (int, error) {
	if _, e := r.conn.Write([]byte(req)); e != nil {
		return 0, e
	}

	m := make([]byte, 1024)
	n, e := r.conn.Read(m)
    if e != nil {
		return 0, e
	}

    resp := string(m[:n])

	f := strings.Fields(resp)
	if len(f) > 2 && strings.HasPrefix(f[0], "RTSP") {
		return strconv.Atoi(f[1])
	}

	return 0, errors.New("Bad response")
}

func (r *RTSP) Query(path string) string {
	method := "DESCRIBE"

	if path == "*" {
		method = "OPTIONS"
	}

	return fmt.Sprintf(_RTSP_TPL, method, path, r.cseq)
}

func (r *RTSP) Check() <-chan string {
	go r.check()
	return r.ch
}

func (r *RTSP) check() {
	defer close(r.ch)
	address := fmt.Sprintf("%s:%d", r.Host, r.Port)
	var err error

	if r.conn, err = net.DialTimeout("tcp", address, r.timeout); err != nil {
		return
	}

	defer r.conn.Close()
	r.conn.SetDeadline(time.Now().Add(time.Second * 5))

	if _, err = r.Request(r.Query("*")); err != nil {
		return
	}

	var code int
	code, err = r.Request(r.Query("/"))
	if err != nil || code == 401 {
		return
	}

	if code == 200 {
		r.ch <- fmt.Sprintf("rtsp://%s/", address)
		return
	}

	for _, path := range r.paths {
		code, err = r.Request(r.Query(path))
		if err != nil || code == 401 {
			return
		}
		if code == 200 {
			r.ch <- fmt.Sprintf("rtsp://%s%s", address, path)
			return
		}
	}
}

func NewRTSP(host string, port int, paths []string, timeout time.Duration) *RTSP {
	return &RTSP{
		Host:    host,
		Port:    port,
		timeout: timeout,
		paths:   paths,
		ch:      make(chan string),
	}
}
