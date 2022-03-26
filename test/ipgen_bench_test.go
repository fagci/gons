package gons_test

import (
	"go-ns/src/generators"
	"testing"
)

func BenchmarkIPGenerator(b *testing.B) {
	g := generators.NewIPGenerator(1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.GenerateWANIP()
	}
}
