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
	Host  string
	Port  int
	conn  net.Conn
	cseq  int
	paths []string
	ch    chan string
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
	if _, e := r.conn.Read(m); e != nil {
		return 0, e
	}

	f := strings.Fields(string(m))
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
	d := net.Dialer{Timeout: time.Second * 2}
	var err error

	if r.conn, err = d.Dial("tcp", address); err != nil {
		return
	}

	defer r.conn.Close()

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

func NewRTSP(host string, port int, paths []string) *RTSP {
	return &RTSP{
		Host:  host,
		Port:  port,
		paths: paths,
		ch:    make(chan string),
	}
}
