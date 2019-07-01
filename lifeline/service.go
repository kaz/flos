package lifeline

import (
	"os/exec"
	"time"

	"github.com/kaz/flos/state"
)

type (
	Result struct {
		Name      string
		Success   bool
		Output    string
		Timestamp int64
	}
)

func runMaster() {
	ch := make(chan *Result)

	for _, cfg := range state.Get().Lifeline {
		logger.Println("script:", cfg.Script)
		go runWorker(cfg.Name, cfg.Script, time.Duration(cfg.Cycle), ch)
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
			Name:      name,
			Success:   err == nil,
			Output:    string(out),
			Timestamp: time.Now().UnixNano(),
		}
		time.Sleep(cycle * time.Second)
	}
}
