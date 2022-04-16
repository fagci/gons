package services

import (
	"bufio"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/fagci/gons/network"
)

type BannerService struct {
	*Service
	timeout time.Duration
}

func (s *BannerService) ScanAddr(addr net.TCPAddr, ch chan<- HostResult, wg *sync.WaitGroup) {
	defer wg.Done()
	if conn, err := network.DialTimeout("tcp", addr.String(), s.timeout); err == nil {
		defer conn.Close()
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		ch <- HostResult{
			Addr: &addr,
			Details: &BannerResult{
				addr:   addr,
				banner: strings.Trim(msg, "\r\n\t "),
			},
		}
	}
}

func NewBannerService(ports []int, timeout time.Duration) *BannerService {
	s := &BannerService{
		timeout: timeout,
		Service: &Service{Ports: ports},
	}
	s.ServiceInterface = interface{}(s).(ServiceInterface)
	return s
}

type BannerResult struct {
	addr   net.TCPAddr
	banner string
}

func (s *BannerResult) String() string {
	return s.addr.String() + "\n" + s.banner
}
func (result *BannerResult) ReplaceVars(cmd string) string {
	return cmd
}
