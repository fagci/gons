package generators

import (
	"fmt"
	"testing"
)

func TestIPGen(t *testing.T) {
	const N = 5
	var i int
	ipGen := NewIPGenerator(1, 5)
	for ip := range ipGen.Generate() {
		fmt.Println(ip.String())
		i++
	}
	if i != N {
		t.Error("Wrong IP count")
	}
}
