package protocol

import (
	"errors"
	"fmt"
	"go-ns/src/generators"
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

const RW_TIMEOUT = time.Second * 10

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

	responseBytes := make([]byte, 1024)
	n, e := r.conn.Read(responseBytes)
	if e != nil {
		return 0, e
	}

	response := string(responseBytes[:n])

	f := strings.Fields(response)
	if len(f) > 2 && strings.HasPrefix(f[0], "RTSP") {
		code, e := strconv.Atoi(f[1])
		/* if e != nil || (code != 200 && code != 401 && code != 404) {
			fmt.Println("Err:", e)
			fmt.Println("Code:", code)
			fmt.Println(req)
			fmt.Println("----")
		} */
		return code, e
	}

	return 0, errors.New("Bad response")
}

func (r *RTSP) Query(path string) string {
	method := "DESCRIBE"

	if path == "*" {
		method = "OPTIONS"
	} else {
		path = fmt.Sprintf("rtsp://%s:%d%s", r.Host, r.Port, path)
	}

	r.cseq++
	return fmt.Sprintf(_RTSP_TPL, method, path, r.cseq)
}

func (r *RTSP) Check() <-chan string {
	r.ch = make(chan string)
	r.cseq = 0
	go r.check()
	return r.ch
}

func (r *RTSP) check() {
	var err error

	defer close(r.ch)

	address := fmt.Sprintf("%s:%d", r.Host, r.Port)

	if r.conn, err = net.DialTimeout("tcp", address, r.timeout); err != nil {
		return
	}

	defer r.conn.Close()

	if _, err = r.Request(r.Query("*")); err != nil {
		return
	}

	var code int
	code, err = r.Request(r.Query(generators.RandomPath(6, 12)))
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
	}
}
