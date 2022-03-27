package protocol

import (
	"errors"
	"fmt"
	"go-ns/src/generators"
	"go-ns/src/models"
	"go-ns/src/utils"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type RTSP struct {
	Host    net.IP
	Port    int
	timeout time.Duration
	conn    net.Conn
	cseq    int
	paths   []string
	ch      chan models.HostResult
}

type RTSPResult struct {
	Url url.URL
}

func (result *RTSPResult) ReplaceVars(cmd string) string {
	cmd = strings.ReplaceAll(cmd, "{result}", result.Url.String())
	cmd = strings.ReplaceAll(cmd, "{scheme}", result.Url.Scheme)
	cmd = strings.ReplaceAll(cmd, "{host}", result.Url.Host)
	cmd = strings.ReplaceAll(cmd, "{hostname}", result.Url.Hostname())
	cmd = strings.ReplaceAll(cmd, "{port}", result.Url.Port())
	cmd = strings.ReplaceAll(cmd, "{slug}", result.Slug())
	return cmd
}

func (result *RTSPResult) Slug() string {
	return utils.Slugify(result.Url.String())
}

func (result *RTSPResult) String() string {
	return result.Url.String()
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

func (r *RTSP) Check() <-chan models.HostResult {
	r.ch = make(chan models.HostResult)
	r.cseq = 0
	go r.check()
	return r.ch
}

func (r *RTSP) result(path string) {
	res := &RTSPResult{}
	res.Url = url.URL{
		Scheme: "rtsp",
		Host:   fmt.Sprintf("%s:%d", r.Host, r.Port),
		Path:   path,
	}
	r.ch <- models.HostResult{
		Addr: &net.TCPAddr{
			IP:   r.Host,
			Port: r.Port,
		},
		Details: res,
	}
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
	if err != nil {
		return
	}

	if code == 200 {
		code, err = r.Request("/")
		if err == nil && code == 200 {
			r.result("/")
		}
		return
	}

	for _, path := range r.paths {
		code, err = r.Request(r.Query(path))
		if err != nil || code == 401 {
			return
		}
		if code == 200 {
			r.result(path)
			return
		}
	}
}

func NewRTSP(host net.IP, port int, paths []string, timeout time.Duration) *RTSP {
	return &RTSP{
		Host:    host,
		Port:    port,
		timeout: timeout,
		paths:   paths,
	}
}
