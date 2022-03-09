package gen

import (
    "net"
    "time"
    "math/rand"
    "encoding/binary"
)

func WanIpGenerator() func() net.IP {
    rand.Seed(time.Now().UnixNano())
    return func() net.IP {
        return RandomWanIP()
    }
}

func RandomWanIP() net.IP {
    for {
        intip := rand.Intn(0xD0000000) + 0xFFFFFF
        if !(
            (intip >= 0x0A000000 && intip <= 0x0AFFFFFF) ||
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
            (intip >= 0xCB007100 && intip <= 0xCB0071FF) ||
            (intip >= 0xF0000000 && intip <= 0xFFFFFFFF)){
            ip := make(net.IP, 4)
            binary.BigEndian.PutUint32(ip, uint32(intip))
            return ip
        }
    }
}
