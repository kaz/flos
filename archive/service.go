package archive

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/kaz/flos/camo"
	"github.com/kaz/flos/libra"

	"github.com/fsnotify/fsnotify"
	"github.com/kaz/flos/state"
	"go.etcd.io/bbolt"
)

type (
	archiver struct {
		watcher *fsnotify.Watcher
		db      *bbolt.DB
	}
)

func runMaster() {
	archiver, err := newArchiver()
	if err != nil {
		logger.Printf("failed to init watcher: %v\n", err)
		return
	}

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

	archiver.Watch("_")
	archiver.Start()
}

func newArchiver() (*archiver, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	db, err := bbolt.Open(DB_FILE, 0644, nil)
	if err != nil {
		logger.Printf("Failed to open db: %v\n", err)
		return nil, err
	}

	return &archiver{watcher, db}, nil
}

func (a *archiver) Start() {
	for ev := range a.watcher.Events {
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
		has, err := a.hasSnapshot(path)
		if err != nil {
			return fmt.Errorf("failed to check snapshot: %v\n", err)
		}
		if !has {
			if err := a.snapshot(path); err != nil {
				return fmt.Errorf("failed to watch children: %v\n", err)
			}
		}
	}
	return nil
}

func (a *archiver) snapshot(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %v\n", err)
	}

	bucketName, err := ptob(path)
	if err != nil {
		return fmt.Errorf("failed to get bucket name: %v\n", err)
	}

	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(time.Now().UnixNano()))

	value, err := camo.Encode(data)
	if err != nil {
		return fmt.Errorf("failed to encode data: %v\n", err)
	}

	return a.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}

		logger.Println("creating snapshot:", path)
		return b.Put(key, value)
	})
}
func (a *archiver) hasSnapshot(path string) (bool, error) {
	bucketName, err := ptob(path)
	if err != nil {
		return false, err
	}

	var has bool
	return has, a.db.View(func(tx *bbolt.Tx) error {
		has = tx.Bucket(bucketName) != nil
		return nil
	})
}

func ptob(p string) ([]byte, error) {
	k, err := camo.Encode([]byte(p))
	if err != nil {
		return nil, fmt.Errorf("failed to encode data: %v\n", err)
	}
	return k, nil
}
