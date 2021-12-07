// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package pkg

import (
	"fmt"

	"github.com/mattn/go-colorable"
)

var (
	// Color and Icon of "CURRENT" output message.
	Icon, Color string

	// Colorable std-out and std-err.
	Stdout = colorable.NewColorableStdout()
	Stderr = colorable.NewColorableStderr()
)

// level is a custom type of `string-level`.
// used define level for [outputLevel] function.S
type level string

// Define constant app levels.
const (
	ErrorL   level = "error"
	SuccessL level = "success"
	InfoL    level = "info"
)

// Define constant color codes.
const (
	RED     string = "\033[0;31m"
	GREEN   string = "\033[0;32m"
	YELLOW  string = "\033[1;33m"
	NOCOLOR string = "\033[0m"
)

// Alert, prints message at given [level].
//
// l - (level) decides style(level) of log message.
// msg - (message) is the content of log message.
func Alert(l level, msg string) {
	// Configure message
	message := fmt.Sprintf("\n %s %s \n", outputLevel(l), msg)

	fmt.Fprintln(Stdout, message)
}

// outputLevel sets [Color] and [Icon] by given `level`,
// and then, returns final printable level title.
//
// Result cases:
// [ERROR] - (powered with red color)
// [SUCCESS] - (powered with green color)
// [INFO] - (powered with yellow color)
func outputLevel(l level) string {
	switch l {
	case ErrorL:
		Color = RED
		Icon = "[ERROR]"
	case SuccessL:
		Color = GREEN
		Icon = "[OK]"
	case InfoL:
		Color = YELLOW
		Icon = "[INFO]"
	default:
		Color = NOCOLOR
	}

	return fmt.Sprintf("%s%s%s", Color, Icon, NOCOLOR)
}
