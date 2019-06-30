package tail

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/kaz/flos/libra"
	"github.com/kaz/flos/state"
)

var (
	logger = log.New(os.Stdout, "[tail] ", log.Ltime)

	// now concurrent read only
	files = map[string]*os.File{}
)

func StartWorker() {
	tailer, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Printf("failed to init watcher: %v\n", err)
		return
	}
	defer tailer.Close()

	s, err := state.RootState().Get("/tail")
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

		file, err := os.Open(path)
		if !ok {
			logger.Printf("failed to open file: %v\n", err)
			continue
		}
		if _, err := file.Seek(0, io.SeekEnd); err != nil {
			logger.Printf("failed to seek file: %v\n", err)
			continue
		}

		files[path] = file

		if err := tailer.Add(path); err != nil {
			logger.Printf("failed to watch: %v\n", err)
			continue
		}
		logger.Printf("Watching file=%v\n", path)
	}

	for ev := range tailer.Events {
		go eventProcess(&ev)
	}
}

func eventProcess(ev *fsnotify.Event) {
	if ev.Op&fsnotify.Write != 0 {
		file, ok := files[ev.Name]
		if !ok {
			logger.Printf("event for unexpected file=%v detected\n", ev.Name)
			return
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			logger.Printf("failed to read file: %v\n", err)
			return
		}
		if _, err := file.Seek(0, io.SeekEnd); err != nil {
			logger.Printf("failed to seek file: %v\n", err)
			return
		}

		libra.Put("tail", fmt.Sprintf("[%v] %v", ev.Name, strings.Trim(string(data), "\r\n")))
	}
}
