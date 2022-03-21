package svc

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type RTSPService struct {
	Port  int
	paths []string
}

type RTSP struct {
	Host  string
	Port  int
	conn  net.Conn
	cseq  int
	paths []string
	ch    chan Result
}

const RTSP_HDR = "%s %s RTSP/1.0\r\n" +
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
	var method string

	if path == "*" {
		method = "OPTIONS"
	} else {
		method = "DESCRIBE"
	}

	return fmt.Sprintf(RTSP_HDR, method, path, r.cseq)
}

func (r *RTSP) check() {
    defer close(r.ch)
	address := fmt.Sprintf("%s:%d", r.Host, r.Port)
	d := net.Dialer{Timeout: time.Second * 2}
	var err error

	r.conn, err = d.Dial("tcp", address)

	if err != nil {
		return
	}

	defer r.conn.Close()

	_, err = r.Request(r.Query("*"))
	if err != nil {
		return
	}

	var code int
	code, err = r.Request(r.Query("/"))
	if err != nil || code == 401 {
		return
	}

	if code == 200 {
		r.ch <- Result{URI: fmt.Sprintf("rtsp://%s/", address)}
		return
	}

	for _, path := range r.paths {
		code, err = r.Request(r.Query(path))
		if err != nil || code == 401 {
			return
		}
		if code == 200 {
			r.ch <- Result{URI: fmt.Sprintf("rtsp://%s%s", address, path)}
			return
		}
	}
}

func (r *RTSP) Check() <-chan Result {
	go r.check()
	return r.ch
}

func (rs *RTSPService) Check(host string) <-chan Result {
    r := NewRTSP(host, rs.Port, rs.paths)
    return r.Check()
}


func NewRTSP(host string, port int, paths []string) *RTSP {
	return &RTSP{
		Host:  host,
		Port:  port,
		paths: paths,
		ch:    make(chan Result),
	}
}

func NewRTSPService(port int, paths []string) *RTSPService {
	return &RTSPService{
		Port:  port,
		paths: paths,
	}
}
