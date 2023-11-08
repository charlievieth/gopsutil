//go:build darwin
// +build darwin

package host

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"unsafe"

	"github.com/shirou/gopsutil/v3/internal/common"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/sys/unix"
)

// from utmpx.h
const user_PROCESS = 7

func HostIDWithContext(ctx context.Context) (string, error) {
	out, err := invoke.CommandWithContext(ctx, "ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	if err != nil {
		return "", err
	}
	if _, line, ok := common.CutBytes(out, []byte(`"IOPlatformUUID" = "`)); ok {
		if j := bytes.IndexByte(line, '"'); j >= 0 {
			return string(line[:j]), nil
		}
	}
	return "", errors.New("cannot find host id")
}

func numProcs(ctx context.Context) (uint64, error) {
	procs, err := process.PidsWithContext(ctx)
	if err != nil {
		return 0, err
	}
	return uint64(len(procs)), nil
}

func UsersWithContext(ctx context.Context) ([]UserStat, error) {
	const utmpfile = "/var/run/utmpx"
	var ret []UserStat

	file, err := os.Open(utmpfile)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return ret, err
	}

	u := Utmpx{}
	const entrySize = int(unsafe.Sizeof(u))
	count := len(buf) / entrySize

	for i := 0; i < count; i++ {
		b := buf[i*entrySize : i*entrySize+entrySize]

		var u Utmpx
		br := bytes.NewReader(b)
		err := binary.Read(br, binary.LittleEndian, &u)
		if err != nil {
			continue
		}
		if u.Type != user_PROCESS {
			continue
		}
		user := UserStat{
			User:     common.IntToString(u.User[:]),
			Terminal: common.IntToString(u.Line[:]),
			Host:     common.IntToString(u.Host[:]),
			Started:  int(u.Tv.Sec),
		}
		ret = append(ret, user)
	}

	return ret, nil
}

var productVersion string
var productVersionOnce sync.Once

func loadProductVersion(ctx context.Context) string {
	productVersionOnce.Do(func() {
		out, err := invoke.CommandWithContext(ctx, "sw_vers", "-productVersion")
		if err == nil {
			productVersion = string(bytes.ToLower(bytes.TrimSpace(out)))
		}
	})
	return productVersion
}

var platformInformation struct {
	platform, family, pver string
	err                    error
	once                   sync.Once
}

func loadPlatformInformationOnce(ctx context.Context) {
	info := &platformInformation
	p, err := unix.Sysctl("kern.ostype")
	if err == nil {
		info.platform = strings.ToLower(p)
	}

	out, err := invoke.CommandWithContext(ctx, "sw_vers", "-productVersion")
	if err == nil {
		info.pver = string(bytes.ToLower(bytes.TrimSpace(out)))
	}

	// check if the macos server version file exists
	_, err = os.Stat("/System/Library/CoreServices/ServerVersion.plist")

	// server file doesn't exist
	if os.IsNotExist(err) {
		info.family = "Standalone Workstation"
	} else {
		info.family = "Server"
	}
}

func PlatformInformationWithContext(ctx context.Context) (string, string, string, error) {
	platformInformation.once.Do(func() {
		loadPlatformInformationOnce(ctx)
	})
	return platformInformation.platform, platformInformation.family,
		platformInformation.pver, platformInformation.err
}

func VirtualizationWithContext(ctx context.Context) (string, string, error) {
	return "", "", common.ErrNotImplementedError
}

func KernelVersionWithContext(ctx context.Context) (string, error) {
	version, err := unix.Sysctl("kern.osrelease")
	return strings.ToLower(version), err
}
