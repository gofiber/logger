package logger

import (
	"testing"
)

// go test -v ./... -run=^$ -bench=Benchmark_Logger -benchmem -count=3

func Benchmark_Logger(b *testing.B) {
	for n := 0; n < b.N; n++ {

	}
}
