package log

import (
	"strings"
)

// [source/uber zap](https://github.com/uber-go/zap/blob/6d482535bdd97f4d97b2f9573ac308f1cf9b574e/zapcore/entry.go#L98C1-L135C1)

// TrimmedPath returns a package/file:line description of the caller,
// preserving only the leaf directory name and file name.
func trimPath(p string, n int) string {
	idx := len(p)
	for i := 0; i < n; i++ {
		idx = strings.LastIndexByte(p[:idx], '/')
	}

	if idx == 0 || idx == -1 {
		return p
	}

	return p[idx+1:]
}
