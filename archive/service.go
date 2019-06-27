package archive

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/kaz/flos/state"
)

type (
	archiver struct {
		*fsnotify.Watcher

		watching map[string][]string
	}
)

func newArchiver() (*archiver, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &archiver{watcher, make(map[string][]string)}, nil
}

func runMaster() {
	archiver, err := newArchiver()
	if err != nil {
		logger.Printf("failed to init watcher: %v\n", err)
		return
	}
	defer archiver.Close()

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

		if err := archiver.Watch(path); err != nil {
			logger.Printf("failed to watch: %v\n", err)
			continue
		}
		logger.Printf("Watching file=%v\n", path)
	}

	archiver.Watch(".")
	archiver.Start()
}

func (a *archiver) Start() {
	for ev := range a.Events {
		logger.Println(ev)

		if ev.Op&fsnotify.Create != 0 {
			if err := a.Watch(ev.Name); err != nil {
				logger.Printf("failed to watch: %v\n", err)
				continue
			}
		}
		if ev.Op&(fsnotify.Remove|fsnotify.Rename) != 0 {
			if err := a.Unwatch(ev.Name); err != nil {
				logger.Printf("failed to unwatch: %v\n", err)
				continue
			}
		}
		if ev.Op&fsnotify.Write != 0 {
			if err := a.snapshot(ev.Name); err != nil {
				logger.Printf("failed to create snapshot: %v\n", err)
				continue
			}
		}
	}
}

func (a *archiver) Watch(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat: %v\n", err)
	}
	return a.watch(info)
}
func (a *archiver) watch(info os.FileInfo) error {
	if info.IsDir() {
		if _, ok := a.watching[info.Name()]; ok {
			return fmt.Errorf("already watching: %v\n", info.Name())
		}

		if err := a.Add(info.Name()); err != nil {
			return fmt.Errorf("failed to add to watcher: %v\n", err)
		}

		children, err := a.watchChildren(info)
		if err != nil {
			return fmt.Errorf("failed to watch children: %v\n", err)
		}
		a.watching[info.Name()] = children
	} else {
		if err := a.snapshot(info.Name()); err != nil {
			return fmt.Errorf("failed to watch children: %v\n", err)
		}
	}
	return nil
}
func (a *archiver) watchChildren(info os.FileInfo) ([]string, error) {
	children := []string{}

	dirs, err := ioutil.ReadDir(info.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %v\n", err)
	}

	for _, ent := range dirs {
		if err := a.watch(ent); err != nil {
			return nil, fmt.Errorf("failed to watch: %v\n", err)
		}

		if ent.IsDir() {
			children = append(children, ent.Name())
		}
	}

	return children, nil
}

func (a *archiver) Unwatch(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat: %v\n", err)
	}
	return a.unwatch(info)
}
func (a *archiver) unwatch(info os.FileInfo) error {
	if info.IsDir() {
		if _, ok := a.watching[info.Name()]; ok {
			return fmt.Errorf("not watching: %v\n", info.Name())
		}

		if err := a.Remove(info.Name()); err != nil {
			return fmt.Errorf("failed to remove from watcher: %v\n", err)
		}

		if err := a.unwatchChildren(a.watching[info.Name()]); err != nil {
			return fmt.Errorf("failed to unwatch children: %v\n", err)
		}

		delete(a.watching, info.Name())
	}
}
func (a *archiver) Unwatch(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat: %v\n", err)
	}
	return a.unwatch(info)
}

func (a *archiver) snapshot(path string) error {
	return nil
}
