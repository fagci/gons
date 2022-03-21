package svc

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type RTSP struct {
	Address string
	conn    net.Conn
	cseq    int
	ch      chan string
}

const PORT = "554"

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

func (r *RTSP) check(paths []string) {
	defer close(r.ch)
	d := net.Dialer{Timeout: time.Second * 2}
	var err error

	r.conn, err = d.Dial("tcp", r.Address)

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
		r.ch <- fmt.Sprintf("rtsp://%s/", r.Address)
		return
	}

	for _, path := range paths {
		code, err = r.Request(r.Query(path))
		if err != nil || code == 401 {
			return
		}
		if code == 200 {
			r.ch <- fmt.Sprintf("rtsp://%s%s", r.Address, path)
			return
		}
	}
}

func (r *RTSP) CheckPaths(paths []string) <-chan string {
	go r.check(paths)
	return r.ch
}

func NewRTSP(address string) *RTSP {
	return &RTSP{
		Address: address,
		ch:      make(chan string),
	}
}
