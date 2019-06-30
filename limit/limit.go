package limit

import (
	"log"
	"os"
	"os/exec"
)

var (
	logger = log.New(os.Stdout, "[limit] ", log.Ltime)
)

func Init() {
	run("ulimit", "-n", "1000000000")
	run("sysctl", "-w", "fs.inotify.max_user_watches=2147483647")
}

func run(name string, arg ...string) {
	out, err := exec.Command(name, arg...).CombinedOutput()
	if err != nil {
		logger.Println(err)
	}
	logger.Println(string(out))
}
