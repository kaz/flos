package archive

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/kaz/flos/state"
)

func runMaster() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Printf("failed to init watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	s, err := state.RootState().Get("/archive")
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

		if err := watcher.Add(path); err != nil {
			logger.Printf("failed to watch: %v\n", err)
			continue
		}
		logger.Printf("Watching file=%v\n", path)
	}

	logger.Println(watcher.Add("."))

	for ev := range watcher.Events {
		fmt.Println(ev)
	}
}
