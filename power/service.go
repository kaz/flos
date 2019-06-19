package power

import (
	"os"
	"os/exec"
)

func stop() {
	logger.Println("Stopping ...")
	os.Exit(0)
}
func restart() {
	me, err := os.Executable()
	if err != nil {
		logger.Printf("failed to get executable name: %v\n", me)
		logger.Printf("fallbacked to os.Args[0] -> %s\n", os.Args[0])
		me = os.Args[0]
	}

	logger.Println("Waiting flesh process ...")
	exec.Command(me).Start()
}
