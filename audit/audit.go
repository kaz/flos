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

	s, err := state.RootState().Get("/audit/file")
	if err != nil {
		logger.Printf("failed to read config: %v\n", err)
		return
	}

	for _, cfg := range s.List() {
		path, ok := cfg.Value().(string)
		if !ok {
			logger.Printf("invalid config type")
			continue
		}
		if err := auditor.WatchFile(path); err != nil {
			logger.Printf("failed to watch: %v\n", err)
			continue
		}
		logger.Printf("Watching file=%v\n", path)
	}

	s, err = state.RootState().Get("/audit/mount")
	if err != nil {
		logger.Printf("failed to read config: %v\n", err)
		return
	}

	for _, cfg := range s.List() {
		path, ok := cfg.Value().(string)
		if !ok {
			logger.Printf("invalid config type")
			return
		}
		if err := auditor.WatchMount(path); err != nil {
			logger.Printf("failed to watch: %v\n", err)
			return
		}
		logger.Printf("Watching mount=%v\n", path)
	}

	for ev := range auditor.Event {
		if bookshelf.IsBookshelf(ev.FileName) {
			continue
		}

		libra.Put("audit", fmt.Sprintln(ev.Acts, ev.FileName, "by", ev.ProcessInfo))
	}
}
