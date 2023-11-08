//go:build linux || freebsd || openbsd || darwin || solaris
// +build linux freebsd openbsd darwin solaris

package host

import (
	"sync"

	"golang.org/x/sys/unix"
)

var cachedKernelArch = struct {
	arch string
	err  error
	once sync.Once
}{}

// type responseCache[T any] struct {
// 	res  T
// 	err  error
// 	once sync.Once
// }

func XKernelArch() (string, error) {
	cachedKernelArch.once.Do(func() {
		var utsname unix.Utsname
		err := unix.Uname(&utsname)
		if err != nil {
			cachedKernelArch.err = err
			return
		}
		cachedKernelArch.arch = unix.ByteSliceToString(utsname.Machine[:])
	})
	return cachedKernelArch.arch, cachedKernelArch.err
}

// TODO: cache this
func KernelArch() (string, error) {
	var utsname unix.Utsname
	err := unix.Uname(&utsname)
	if err != nil {
		return "", err
	}
	return unix.ByteSliceToString(utsname.Machine[:]), nil
}
