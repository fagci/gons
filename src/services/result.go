package services

import (
	"fmt"
	"github.com/fagci/gons/src/utils"
	"net"
	"strings"
)

type ResultDetails interface {
	ReplaceVars(string) string
	String() string
}

type HostResult struct {
	Addr    net.Addr
	Details ResultDetails
}

func (result *HostResult) ReplaceVars(cmd string) string {
	if result.Details != nil {
		cmd = result.Details.ReplaceVars(cmd)
	}
	host := result.Addr.String()
	var hostname string
	var port int
	switch addr := result.Addr.(type) {
	case *net.TCPAddr:
		hostname = addr.IP.String()
		port = addr.Port
	case *net.UDPAddr:
		hostname = addr.IP.String()
		port = addr.Port
	}
	if port == 0 {
		host = hostname
	}
	cmd = strings.ReplaceAll(cmd, "{hostname}", hostname)
	cmd = strings.ReplaceAll(cmd, "{host}", host)
	cmd = strings.ReplaceAll(cmd, "{port}", fmt.Sprintf("%d", port))
	cmd = strings.ReplaceAll(cmd, "{proto}", result.Addr.Network())
	return cmd
}

func (result *HostResult) String() string {
	if result.Details != nil {
		return result.Details.String()
	}
	switch addr := result.Addr.(type) {
	case *net.UDPAddr:
		if addr.Port == 0 {
			return addr.IP.String()
		}
	case *net.TCPAddr:
		if addr.Port == 0 {
			return addr.IP.String()
		}
	}
	return result.Addr.String()
}

func (result *HostResult) Slug() string {
	return utils.Slugify(result.Addr.String())
}
