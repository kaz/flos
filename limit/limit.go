package limit

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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
	out, err := exec.Command("bash", "-c", `
		min=0
		max=$(( 1 << 62 ))

		while (( max - min > 1 )); do
			mid=$(( (min + max) >> 1 ))
			if echo $mid > /proc/sys/fs/inotify/max_user_watches; then
				min=$mid
			else
				max=$mid
			fi
		done

		echo $mid | tee /proc/sys/fs/inotify/max_user_watches
	`).Output()
	if err != nil {
		return err
	}

	logger.Printf("fs.inotify.max_user_watches set to %s\n", out)
	return nil
}
