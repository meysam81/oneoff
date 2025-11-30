//go:build windows

package jobs

import (
	"os/exec"

	"github.com/meysam81/oneoff/internal/logging"
)

func setSysProcAttr(cmd *exec.Cmd) {
}

func killProcessGroup(cmd *exec.Cmd) {
	if cmd.Process != nil {
		err := cmd.Process.Kill()
		logging.Error().Err(err).Msg("failed to kill the process")
	}
}
