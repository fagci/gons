package network

import (
	"context"
	"net"
	"time"
)

var dialService = struct {
	LocalAddr *net.TCPAddr
}{}

func SetInterface(iface string) error {
	if iface != "" {
		ief, err := net.InterfaceByName(iface)
		if err != nil {
			return err
		}

		addrs, err := ief.Addrs()
		if err != nil {
			return err
		}

		addr := &net.TCPAddr{
			IP: addrs[0].(*net.IPNet).IP,
		}

		dialService.LocalAddr = addr
	}

	return nil
}

func DialTimeout(network string, address string, timeout time.Duration) (net.Conn, error) {
	_d := net.Dialer{
		Timeout:   timeout,
		LocalAddr: dialService.LocalAddr,
	}

	return _d.Dial(network, address)
}

func DialContextFunc(timeout time.Duration) func(ctx context.Context, network, address string) (net.Conn, error) {
	_d := net.Dialer{
		Timeout:   timeout,
		LocalAddr: dialService.LocalAddr,
	}

	return _d.DialContext
}
