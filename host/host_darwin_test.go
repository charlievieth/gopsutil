package host

import (
	"context"
	"testing"
)

func BenchmarkHostIDWithContext(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, err := HostIDWithContext(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkKernelVersionWithContext(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, err := KernelVersionWithContext(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}
