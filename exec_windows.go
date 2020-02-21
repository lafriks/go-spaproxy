// +build windows

package spaproxy

import (
	"os/exec"
)

func newCommand(path string, args ...string) *exec.Command {
	a := make([]string, 0, len(args)+2)
	a = append(a, "/c", path)
	a = append(a, args...)
	return exec.NewCommand("cmd.exe", a...)
}
