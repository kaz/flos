package archive

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/kaz/flos/libra"
	"github.com/kaz/flos/libra/bookshelf"
)

const (
	ARCHIVE_FILE  = "chunk.0003.zip"
	MAX_ROW_COUNT = 1 << 10
)

type (
	archiver struct {
		watcher *fsnotify.Watcher
		shelf   *bookshelf.Bookshelf
	}
)

func NewArchiver() (*archiver, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	shelf, err := bookshelf.New(ARCHIVE_FILE, MAX_ROW_COUNT)
	if err != nil {
		logger.Printf("Failed to open db: %v\n", err)
		return nil, err
	}

	return &archiver{watcher, shelf}, nil
}

func (a *archiver) Start() {
	for ev := range a.watcher.Events {
		if bookshelf.IsBookshelf(ev.Name) {
			continue
		}

		libra.Put("archive", ev.String())

		if ev.Op&fsnotify.Create != 0 {
			if err := a.Watch(ev.Name); err != nil {
				logger.Printf("failed to watch: %v\n", err)
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
	abspath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get path: %v\n", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat: %v\n", err)
	}

	return a.watch(abspath, info)
}
func (a *archiver) watch(path string, info os.FileInfo) error {
	if info.IsDir() {
		if err := a.watcher.Add(path); err != nil {
			return fmt.Errorf("failed to add to watcher: %v\n", err)
		}

		dirs, err := ioutil.ReadDir(path)
		if err != nil {
			return fmt.Errorf("failed to read dir: %v\n", err)
		}

		for _, ent := range dirs {
			if err := a.watch(filepath.Join(path, ent.Name()), ent); err != nil {
				return fmt.Errorf("failed to watch child: %v\n", err)
			}
		}
	} else {
		if bookshelf.IsBookshelf(path) {
			return nil
		} else if has, err := a.hasSnapshot(path); err != nil {
			return fmt.Errorf("failed to check snapshot: %v\n", err)
		} else if has {
			return nil
		} else if err := a.snapshot(path); err != nil {
			return fmt.Errorf("failed to watch children: %v\n", err)
		}
	}
	return nil
}

func (a *archiver) snapshot(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %v\n", err)
	}

	return a.shelf.Put([]byte(path), data)
}
func (a *archiver) hasSnapshot(path string) (bool, error) {
	return a.shelf.Has([]byte(path))
}
