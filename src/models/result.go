package models

import (
	"gons/src/utils"
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
	switch addr := result.Addr.(type) {
	case *net.UDPAddr:
		if addr.Port == 0 {
			host = addr.IP.String()
		}
	case *net.TCPAddr:
		if addr.Port == 0 {
			host = addr.IP.String()
		}
	}
	cmd = strings.ReplaceAll(cmd, "{host}", host)
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
