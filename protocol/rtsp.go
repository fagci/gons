package protocol

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fagci/gons/network"
)

type RTSP struct {
	timeout  time.Duration
	conn     net.Conn
	cseq     int
	paths    []string
	addr     string
	fakePath string
}

const _RTSP_TPL = "%s %s RTSP/1.0\r\n" +
	"Accept: application/sdp\r\n" +
	"CSeq: %d\r\n" +
	"User-Agent: Lavf59.16.100\r\n\r\n"

func (r *RTSP) Request(req string) (cose int, err error) {
	_ = r.conn.SetDeadline(time.Now().Add(RW_TIMEOUT))

	if _, err = r.conn.Write([]byte(req)); err != nil {
		return
	}

	data := make([]byte, 1024)
	n, err := r.conn.Read(data)
	if err != nil {
		return
	}

	f := strings.Fields(string(data[:n]))
	if len(f) >= 2 && strings.HasPrefix(f[0], "RTSP/") {
		return strconv.Atoi(f[1])
	}

	return 0, errors.New("Bad response")
}

func (r *RTSP) Query(path string) string {
	r.cseq++

	return fmt.Sprintf(_RTSP_TPL, "DESCRIBE", fmt.Sprintf("rtsp://%s%s", r.addr, path), r.cseq)
}

func (r *RTSP) Check() (url url.URL, err error) {
	var code int

	if r.conn, err = network.DialTimeout("tcp", r.addr, r.timeout); err != nil {
		return
	}

	defer r.conn.Close()

	code, err = r.Request(r.Query(r.fakePath))
	if err != nil {
		return
	}

	if code == 200 {
		return url, errors.New("RTSP: fake")
	}

	for _, path := range r.paths {
		code, err = r.Request(r.Query(path))
		if err != nil {
			return
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

func NewRTSP(addr net.Addr, paths []string, fakePath string, timeout time.Duration) *RTSP {
	return &RTSP{
		timeout:  timeout,
		paths:    paths,
		addr:     addr.String(),
		fakePath: fakePath,
	}
}
