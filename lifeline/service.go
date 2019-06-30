package lifeline

import (
	"os/exec"
	"time"

	"github.com/kaz/flos/state"
)

type (
	Result struct {
		Timestamp time.Time
		Name      string
		Success   bool
		Output    string
	}
)

func runMaster() {
	s, err := state.RootState().Get("/lifeline")
	if err != nil {
		logger.Printf("failed to read config: %v\n", err)
		return
	}

	ch := make(chan *Result)

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

		logger.Println("script:", script)
		go runWorker(name, script, time.Duration(cycle), ch)
	}

	for r := range ch {
		go resultProcess(r)
	}
}

func resultProcess(r *Result) {
	mu.Lock()
	defer mu.Unlock()

	results[r.Name] = r
}

func runWorker(name, script string, cycle time.Duration, ch chan *Result) {
	for {
		out, err := exec.Command("sh", "-c", script).CombinedOutput()
		if err != nil {
			// logger.Printf("command failed: %v\n", err)
		}

		ch <- &Result{
			Timestamp: time.Now(),
			Name:      name,
			Success:   err == nil,
			Output:    string(out),
		}
		time.Sleep(cycle * time.Second)
	}
}
