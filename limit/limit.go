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
	logger.Println(exec.Command("ulimit", "-n", "1000000000").CombinedOutput())
	logger.Println(exec.Command("sysctl", "-w", "fs.inotify.max_user_watches=2147483647").CombinedOutput())
}
