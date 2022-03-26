package generators

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"net"
)

type IPGenerator struct {
	ch chan net.IP
	r  *rand.Rand
}

func (g *IPGenerator) GenerateWANIP() net.IP {
	var intip uint32
	for {
		intip = g.r.Uint32()%0xD0000000 + 0xFFFFFF
		if (intip >= 0x0A000000 && intip <= 0x0AFFFFFF) ||
			(intip >= 0x64400000 && intip <= 0x647FFFFF) ||
			(intip >= 0x7F000000 && intip <= 0x7FFFFFFF) ||
			(intip >= 0xA9FE0000 && intip <= 0xA9FEFFFF) ||
			(intip >= 0xAC100000 && intip <= 0xAC1FFFFF) ||
			(intip >= 0xC0000000 && intip <= 0xC0000007) ||
			(intip >= 0xC00000AA && intip <= 0xC00000AB) ||
			(intip >= 0xC0000200 && intip <= 0xC00002FF) ||
			(intip >= 0xC0A80000 && intip <= 0xC0A8FFFF) ||
			(intip >= 0xC6120000 && intip <= 0xC613FFFF) ||
			(intip >= 0xC6336400 && intip <= 0xC63364FF) ||
			(intip >= 0xCB007100 && intip <= 0xCB0071FF) {
			continue
		}
		break
	}
	return net.IPv4(byte(intip>>24), byte(intip>>16), byte(intip>>8), byte(intip))
}

func (g *IPGenerator) GenerateWAN() <-chan net.IP {
	go func() {
		g.ch <- g.GenerateWANIP()
	}()

	return g.ch
}
func (g *IPGenerator) Stop() {
	close(g.ch)
}

func NewIPGenerator(capacity int) *IPGenerator {
	b := make([]byte, 8)
	_, err := crypto_rand.Read(b)
	if err != nil {
        panic("Cryptorandom seed failed: " + err.Error())
	}
	return &IPGenerator{
		ch: make(chan net.IP, capacity),
		r:  rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b)))),
	}
}
