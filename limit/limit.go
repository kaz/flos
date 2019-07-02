// +build linux darwin

package limit

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"syscall"
)

var (
	logger = log.New(os.Stdout, "[limit] ", log.Ltime)
)

func Init() {
	if err := setNoFileLimit(); err != nil {
		logger.Printf("failed to set RLIMIT_NOFILE: %v\n", err)
	}
	if err := setInotifyLimit(); err != nil {
		logger.Printf("failed to set fs.inotify.max_user_watches: %v\n", err)
	}
}

func setNoFileLimit() error {
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		return fmt.Errorf("failed to get rlimit: %v\n", err)
	}

	min := limit.Cur
	max := limit.Max

	for max-min > 1 {
		mid := (min + max) / 2
		limit.Cur = mid

		syscall.Setrlimit(syscall.RLIMIT_NOFILE, &limit)
		if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
			return fmt.Errorf("failed to get rlimit: %v\n", err)
		}

		if limit.Cur == mid {
			min = mid
		} else {
			max = mid
		}
	}

	logger.Printf("RLIMIT_NOFILE set to %v\n", limit)
	return nil
}

func setInotifyLimit() error {
	var min uint64 = 0
	var max uint64 = 1<<64 - 1

	for max-min > 1 {
		mid := (min + max) / 2

		if err := ioutil.WriteFile("/proc/sys/fs/inotify/max_user_watches", []byte(strconv.FormatUint(mid, 10)), 0644); err == nil {
			min = mid
		} else {
			max = mid
		}
	}

	logger.Printf("fs.inotify.max_user_watches set to %d\n", min)
	return nil
}
