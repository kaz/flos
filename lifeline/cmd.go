// +build !windows

package lifeline

import "os/exec"

func command(script string) *exec.Cmd {
	return exec.Command("sh", "-c", script)
}
