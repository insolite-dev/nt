// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package pkg

import (
	"fmt"

	"github.com/anonistas/notya/lib/models"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

var (
	// Color and Icon of logger's "CURRENT" output message.
	Icon, Color string

	// ColorableStd is main stdargs of logger. [colorable stdargs].
	ColorableStd = models.StdArgs{
		Stdout: colorable.NewColorableStdout(),
		Stderr: colorable.NewColorableStderr(),
	}
)

// Level is a custom type of [string-level].
// used to define level for [OutputLevel] function.
type Level string

// Defined constant app Levels.
const (
	ErrorL   Level = "error"
	SuccessL Level = "success"
	InfoL    Level = "info"
)

// Defined constant color codes, for [OutputLevel].
const (
	RED     string = "\033[0;31m"
	GREEN   string = "\033[0;32m"
	YELLOW  string = "\033[1;33m"
	NOCOLOR string = "\033[0m"
)

// Defined constant icon/title codes.
const (
	ERROR   string = "[ERROR]"
	SUCCESS string = "[OK]"
	INFO    string = "[INFO]"
)

// Loggers powered by colors.
var (
	divider     = color.New(color.FgHiYellow)
	text        = color.New(color.FgHiWhite)
	lowText     = color.New(color.Faint)
	rainbowText = color.New(color.FgHiMagenta)
)

// Alert, logs message at given [Level].
//
// l - (Level) decides style(Level) of log message.
// msg - (message) is the content of log message.
func Alert(l Level, msg string) {
	// Configure message.
	message := fmt.Sprintf("\n %s %s \n", OutputLevel(l), msg)

	fmt.Fprintln(ColorableStd.Stdout, message)
}

// OutputLevel sets [Color] and [Icon] by given [Level],
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

// ShowNote, logs given full note.
func ShowNote(note models.Note) {
	// Modify note fields to make it ready to log.
	title := fmt.Sprintf("\nTitle: %v", note.Title)
	path := fmt.Sprintf("Path: %v", note.Path)
	body := fmt.Sprintf("\n%v", note.Body)

	// Log the final note files.
	rainbowText.Println(title)
	lowText.Println(path)
	divider.Println("----------------------")

	// Printout no content if body is empty.
	if len(note.Body) == 0 {
		text.Println("No content ...")
	} else {
		text.Println(body)
	}

	divider.Println("----------------------")
}

// ShowListOfNotes, logs given list as limited-element-per-row array.
func ShowListOfNotes(list []string, limit int) {
	if len(list) == 0 {
		return
	}

	var l string

	nl := limit
	for i, e := range list {
		// Limit row by [nl] element.
		if i >= nl {
			l += "\n"
			nl += limit
		}

		l += fmt.Sprintf(" %v ", e)
	}

	div := "----------------------"

	divider.Println(div + div)
	text.Println(l)
	divider.Println(div + div)
}
