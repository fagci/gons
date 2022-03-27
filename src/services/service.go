package services

import (
	"go-ns/src/models"
	"net"
)

type Service interface {
    Check(net.IP) <-chan models.HostResult
}
