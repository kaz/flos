package audit

import (
	"fmt"
	"log"
	"os"

	"github.com/kaz/flos/libra"
	"github.com/kaz/flos/libra/bookshelf"
	"github.com/kaz/flos/state"
)

var (
	logger = log.New(os.Stdout, "[audit] ", log.Ltime)
)

func StartWorker() {
	auditor, err := NewAuditor(false)
	if err != nil {
		logger.Printf("failed to init auditor: %v\n", err)
		return
	}

	for _, path := range state.Get().Audit.File {
		if err := auditor.WatchFile(path); err != nil {
			logger.Printf("failed to watch: %v\n", err)
			continue
		}
		logger.Printf("Watching file=%v\n", path)
	}
	for _, path := range state.Get().Audit.Mount {
		if err := auditor.WatchMount(path); err != nil {
			logger.Printf("failed to watch: %v\n", err)
			return
		}
		logger.Printf("Watching mount=%v\n", path)
	}

	for ev := range auditor.Event {
		go eventProcess(ev)
	}
}

func eventProcess(ev *Event) {
	if bookshelf.IsBookshelf(ev.FileName) {
		return
	}
	libra.Put("audit", fmt.Sprintln(ev.Acts, ev.FileName, "by", ev.ProcessInfo))
}
