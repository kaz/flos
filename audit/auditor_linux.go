package audit

import (
	"fmt"
	"os"

	"github.com/jandre/fanotify"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/unix"
)

type (
	Auditor struct {
		perm  bool
		nd    *fanotify.NotifyFD
		Event chan *Event
	}
	Event struct {
		Acts        []string
		FileName    string
		ProcessInfo string
	}
)

/*
	!!! IMPORTANT !!!
	Sometimes perm=true unexpectedly makes whole system hanged-up.
	I do not realize why the problem occurs, so please be careful to use it :P
*/
func NewAuditor(perm bool) (*Auditor, error) {
	flag := fanotify.FAN_CLASS_NOTIF
	if perm {
		flag = fanotify.FAN_CLASS_PRE_CONTENT
	}

	nd, err := fanotify.Initialize(flag|fanotify.FAN_CLOEXEC|fanotify.FAN_UNLIMITED_QUEUE|fanotify.FAN_UNLIMITED_MARKS, unix.O_RDONLY|unix.O_LARGEFILE)
	if err != nil {
		return nil, err
	}

	a := &Auditor{
		perm:  perm,
		nd:    nd,
		Event: make(chan *Event),
	}

	go a.startAudit()
	return a, nil
}

func (a *Auditor) watch(path string, addFlag int) error {
	evMask := fanotify.FAN_ALL_EVENTS
	if a.perm {
		evMask |= fanotify.FAN_ALL_PERM_EVENTS
	}
	return a.nd.Mark(fanotify.FAN_MARK_ADD|addFlag, uint64(evMask), unix.AT_FDCWD, path)
}
func (a *Auditor) WatchFile(path string) error {
	return a.watch(path, 0)
}
func (a *Auditor) WatchMount(path string) error {
	return a.watch(path, fanotify.FAN_MARK_MOUNT)
}

func (a *Auditor) startAudit() {
	for {
		ev, err := a.nd.GetEvent()
		if err != nil {
			logger.Println(err)
			continue
		}

		var procInfo string
		process, err := process.NewProcess(ev.Pid)
		if err != nil {
			// logger.Println(err)
			procInfo = "[unknown process]"
		} else {
			procInfo = getProcInfo(process)
		}

		fileName, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", ev.File.Fd()))
		if err != nil {
			// logger.Println(err)
			fileName = "[unknown file]"
		}

		acts := []string{}
		if ev.Mask&fanotify.FAN_ACCESS != 0 {
			acts = append(acts, "ACCESS")
		}
		if ev.Mask&fanotify.FAN_OPEN != 0 {
			acts = append(acts, "OPEN")
		}
		if ev.Mask&fanotify.FAN_MODIFY != 0 {
			acts = append(acts, "MODIFY")
		}
		if ev.Mask&fanotify.FAN_CLOSE_WRITE != 0 {
			acts = append(acts, "CLOSE_WRITE")
		}
		if ev.Mask&fanotify.FAN_CLOSE_NOWRITE != 0 {
			acts = append(acts, "CLOSE_NOWRITE")
		}
		if ev.Mask&fanotify.FAN_Q_OVERFLOW != 0 {
			acts = append(acts, "Q_OVERFLOW")
		}
		if ev.Mask&fanotify.FAN_OPEN_PERM != 0 {
			acts = append(acts, "OPEN_PERM")
			a.nd.Response(ev, true)
		}
		if ev.Mask&fanotify.FAN_ACCESS_PERM != 0 {
			acts = append(acts, "ACCESS_PERM")
			a.nd.Response(ev, true)
		}

		ev.File.Close()
		a.Event <- &Event{acts, fileName, procInfo}
	}
}

func getProcInfo(p *process.Process) string {
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

	return info
}
