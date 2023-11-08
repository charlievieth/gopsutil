package common

import (
	"bytes"
	"strings"
)

// TODO: remove any if unused and consider calling the stdlib funcs.

func CutString(s, sep string) (before, after string, found bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

func CutPrefixString(s, prefix string) (after string, found bool) {
	if !strings.HasPrefix(s, prefix) {
		return s, false
	}
	return s[len(prefix):], true
}

func CutSuffixString(s, suffix string) (before string, found bool) {
	if !strings.HasSuffix(s, suffix) {
		return s, false
	}
	return s[:len(s)-len(suffix)], true
}

func CutBytes(s, sep []byte) (before, after []byte, found bool) {
	if i := bytes.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, nil, false
}

func CutPrefixBytes(s, prefix []byte) (after []byte, found bool) {
	if !bytes.HasPrefix(s, prefix) {
		return s, false
	}
	return s[len(prefix):], true
}

func CutSuffixBytes(s, suffix []byte) (before []byte, found bool) {
	if !bytes.HasSuffix(s, suffix) {
		return s, false
	}
	return s[:len(s)-len(suffix)], true
}
