package lifeline

import (
	"os/exec"
	"time"

	"github.com/kaz/flos/state"
)

type (
	result struct {
		timestamp time.Time
		name      string
		result    bool
	}
)

func runMaster() {
	s, err := state.RootState().Get("/lifetime")
	if err != nil {
		logger.Printf("failed to read config: %v\n", err)
		return
	}

	ch := make(chan *result)

	for _, cfg := range s.List() {
		naState, err := cfg.Get("/name")
		if err != nil {
			logger.Printf("failed to read config: %v\n", err)
			continue
		}

		name, ok := naState.Value().(string)
		if !ok {
			logger.Printf("invalid config type")
			continue
		}

		scState, err := cfg.Get("/script")
		if err != nil {
			logger.Printf("failed to read config: %v\n", err)
			continue
		}

		script, ok := scState.Value().(string)
		if !ok {
			logger.Printf("invalid config type")
			continue
		}

		cyState, err := cfg.Get("/cycle")
		if err != nil {
			logger.Printf("failed to read config: %v\n", err)
			continue
		}

		cycle, ok := cyState.Value().(float64)
		if !ok {
			logger.Printf("invalid config type")
			continue
		}

		go runWorker(name, script, time.Duration(cycle), ch)
	}

	for r := range ch {
		status[r.name] = r
	}
}

func runWorker(name, script string, cycle time.Duration, ch chan *result) {
	for {
		err := exec.Command("sh", "-c", script).Run()
		if err != nil {
			logger.Printf("command failed: %v\n", err)
		}

		ch <- &result{time.Now(), name, err == nil}
		time.Sleep(cycle * time.Second)
	}
}
