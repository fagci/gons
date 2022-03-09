package main

import (
	"fmt"
	"go_ns/src/gen"
	"time"
)

func main() {
	ip_gen := gen.WanIpGenerator()

	N := 10000000
    var start time.Time

	for t := 1; t <= 5; t++ {
		start = time.Now()
		for i := 0; i < N; i++ {
			ip_gen()
		}
		fmt.Printf("%d. %d ms\n", t, time.Since(start).Milliseconds())
	}
}
