// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package pkg

import (
	"fmt"

	"github.com/anonistas/notya/lib/models"
	"github.com/mattn/go-colorable"
)

var (
	// Color and Icon of "CURRENT" output message.
	Icon, Color string

	// ColorableStd stores colorable std-out and std-err.
	ColorableStd = models.StdArgs{
		Stdout: colorable.NewColorableStdout(),
		Stderr: colorable.NewColorableStderr(),
	}
)

// Level is a custom type of `string-level`.
// used define level for [OutputLevel] function.
type Level string

// Define constant app Levels.
const (
	ErrorL   Level = "error"
	SuccessL Level = "success"
	InfoL    Level = "info"
)

// Define constant color codes.
const (
	RED     string = "\033[0;31m"
	GREEN   string = "\033[0;32m"
	YELLOW  string = "\033[1;33m"
	NOCOLOR string = "\033[0m"
)

// Define constant icon/title codes.
const (
	ERROR   string = "[ERROR]"
	SUCCESS string = "[OK]"
	INFO    string = "[INFO]"
)

// Alert, prints message at given [Level].
//
// l - (Level) decides style(Level) of log message.
// msg - (message) is the content of log message.
func Alert(l Level, msg string) {
	// Configure message
	message := fmt.Sprintf("\n %s %s \n", OutputLevel(l), msg)

	fmt.Fprintln(ColorableStd.Stdout, message)
}

// OutputLevel sets [Color] and [Icon] by given `Level`,
// and then, returns final printable Level title.
//
// Result cases:
// [ERROR] - (powered with red color)
// [OK] - (powered with green color)
// [INFO] - (powered with yellow color)
func OutputLevel(l Level) string {
	switch l {
	case ErrorL:
		Color = RED
		Icon = ERROR
	case SuccessL:
		Color = GREEN
		Icon = SUCCESS
	case InfoL:
		Color = YELLOW
		Icon = INFO
	default:
		Color = NOCOLOR
		Icon = ""
	}

	return fmt.Sprintf("%s%s%s", Color, Icon, NOCOLOR)
}
