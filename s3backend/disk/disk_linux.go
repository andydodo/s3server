// +build linux

package s3disk

import (
	"os"
	"syscall"
	"time"
)

func GetCTime(fi os.FileInfo) time.Time {
	stat := fi.Sys().(*syscall.Stat_t)

	return time.Unix(stat.Ctime.Unix())
}
