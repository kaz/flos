package archive

import (
	"log"
	"os"

	"github.com/kaz/flos/state"
	"github.com/labstack/echo/v4"
)

var (
	logger = log.New(os.Stdout, "[archive] ", log.Ltime)
)

func StartService(g *echo.Group) {
	archiver, err := NewArchiver()
	if err != nil {
		logger.Printf("failed to init watcher: %v\n", err)
		return
	}

	g.PATCH("/snapshots", archiver.shelf.GetHandler)
	g.DELETE("/snapshots", archiver.shelf.DeleteHandler)

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
		logger.Printf("Watching dir=%v\n", path)
	}

	go archiver.Start()
}
