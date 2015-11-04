// +build darwin freebsd

package s3disk

import (
	"os"
	"syscall"
	"time"
)

func GetCTime(fi os.FileInfo) time.Time {
	stat := fi.Sys().(*syscall.Stat_t)

	return time.Unix(stat.Ctimespec.Unix())
}
