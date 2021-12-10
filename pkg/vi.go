// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package pkg

import (
	"os/exec"

	"github.com/anonistas/notya/lib/models"
)

// OpenWithVI opens file from given path with vi/vim.
func OpenFileWithVI(filepath string, stdargs models.StdArgs) error {
	// Look VI execution path from current running machine.
	vi, pathErr := exec.LookPath("vi")
	if pathErr != nil {
		return pathErr
	}

	// Generate vi command to open file.
	viCmd := &exec.Cmd{
		Path:   vi,
		Args:   []string{vi, filepath},
		Stdin:  stdargs.Stdin,
		Stdout: stdargs.Stdout,
		Stderr: stdargs.Stderr,
	}

	if err := viCmd.Run(); err != nil {
		return err
	}

	return nil
}
