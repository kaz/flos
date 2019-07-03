package audit

import (
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/process"
)

var (
	mu    = sync.RWMutex{}
	cache = map[int32]string{}
)

func cacheFlusher() {
	for {
		mu.Lock()
		for k, _ := range cache {
			delete(cache, k)
		}
		mu.Unlock()

		time.Sleep(2 * time.Minute)
	}
}

func GetProcInfo(pid int32) (string, error) {
	mu.RLock()
	if info, ok := cache[pid]; ok {
		return info, nil
	}
	mu.RUnlock()

	info, err := procInfo(pid)
	if err != nil {
		return "", err
	}

	mu.Lock()
	cache[pid] = info
	mu.Unlock()

	return info, nil
}

func procInfo(pid int32) (string, error) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return "", err
	}

	info := ""

	name, err := p.Name()
	if err == nil {
		info += fmt.Sprintf("%s pid=%d ", name, p.Pid)
	} else {
		info += fmt.Sprintf("(unrecognized) pid=%d ", p.Pid)
	}

	uids, err := p.Uids()
	if err == nil {
		info += fmt.Sprintf("uid=%v ", uids[0])
	} else {
		info += "uid=? "
	}

	gids, err := p.Gids()
	if err == nil {
		info += fmt.Sprintf("gid=%v ", gids[0])
	} else {
		info += "gid=? "
	}

	exe, err := p.Exe()
	if err == nil {
		info += fmt.Sprintf("bin=%s ", exe)
	} else {
		info += "bin=? "
	}

	return info, err
}
