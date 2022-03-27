package services

import (
	"gons/src/models"
	"net"
)

type Service interface {
    Check(net.IP) <-chan models.HostResult
}
