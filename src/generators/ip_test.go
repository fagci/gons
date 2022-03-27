package generators

import (
	"fmt"
	"testing"
)

func TestIPGen(t *testing.T) {
    fmt.Println("Start")
    ipGen := NewIPGenerator(1, 5)
    for ip := range ipGen.GenerateWAN() {
        fmt.Println(ip.String())
    }
}
