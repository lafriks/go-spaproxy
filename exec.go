// +build !windows

package spaproxy

import (
	"context"
	"os/exec"
)

func newCommand(ctx context.Context, path string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, path, args...)
}
