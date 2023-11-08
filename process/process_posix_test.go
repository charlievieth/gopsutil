//go:build linux || freebsd || darwin
// +build linux freebsd darwin

package process

import (
	"os"
	"testing"

	"golang.org/x/sys/unix"
)

func Test_SendSignal(t *testing.T) {
	checkPid := os.Getpid()

	p, _ := NewProcess(int32(checkPid))
	err := p.SendSignal(unix.SIGCONT)
	if err != nil {
		t.Errorf("send signal %v", err)
	}
}

func BenchmarkIsMount(b *testing.B) {
	pwd, err := os.Getwd()
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		// isMount(".")
		isMount(pwd)
	}
}
