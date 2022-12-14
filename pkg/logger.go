//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package pkg

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/insolite-dev/notya/lib/models"
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
	GREY       string = "\033[1;30m"
	RED        string = "\033[0;31m"
	GREEN      string = "\033[0;32m"
	YELLOW     string = "\033[1;33m"
	DARKYELLOW string = "\033[2;33m"
	PURPLE     string = "\033[1;35m"
	CYAN       string = "\033[1;36m"
	NOCOLOR    string = "\033[0m"
)

// Defined constant icon/title codes.
const (
	ERROR   string = "[X]"
	SUCCESS string = "[✔]"
	INFO    string = "[I]"
)

// Loggers powered by colors.
var (
	// divider     = color.New(color.FgHiYellow)
	text = color.New(color.FgHiWhite)
	// lowText     = color.New(color.Faint)
	// rainbowText = color.New(color.FgHiMagenta)
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

// Print, prints given data by combining it with given color attribute.
func Print(data string, c color.Attribute) {
	color.New(c).Println(data)
}

// ShowNote, logs given full note.
func PrintNote(note models.Note) {
	// Modify note fields to make it ready to log.
	title := fmt.Sprintf(
		"\n%v %v",
		fmt.Sprintf("%s%s%s", PURPLE, "Title:", NOCOLOR),
		fmt.Sprintf("%s%s%s", GREY, note.Title, NOCOLOR),
	)
	path := fmt.Sprintf("%v %v",
		fmt.Sprintf("%s%s%s", PURPLE, "Path:", NOCOLOR),
		fmt.Sprintf("%s%s%s", GREY, note.Path, NOCOLOR),
	)

	body := fmt.Sprintf("\n%v", note.Body)

	// Log the final note files.
	text.Println(title)
	if len(note.Path) > 0 {
		text.Println(path)
	}

	// Printout no content if body is empty.
	if len(note.Body) == 0 {
		text.Add(color.FgHiYellow).Println("\n No content ... \n ")
	} else {
		text.Println(body)
	}
}

// PrintNotes, logs given nodes list.
func PrintNodes(list []models.Node) {
	if len(list) == 0 {
		return
	}

	for _, value := range list {
		note := fmt.Sprintf(
			" %v %s %v",
			fmt.Sprintf("%s%s%s", GREY, "•", NOCOLOR),
			fmt.Sprintf("%s%s%s", YELLOW, value.Pretty[0], NOCOLOR),
			fmt.Sprintf("%s%s%s", DARKYELLOW, value.Pretty[1], NOCOLOR),
		)
		text.Println(note)
	}
}

// PrintSettings, logs given settings model.
func PrintSettings(settings models.Settings) {
	values := settings.ToJSON()

	for key, value := range values {
		printable := fmt.Sprintf(" • %s: %s", fmt.Sprintf("%s%s%s", YELLOW, key, NOCOLOR), value)
		text.Println(printable)
	}
}

// PrintErrors, is general error logger for push and fetch command error results.
func PrintErrors(act string, errs []error) {
	for i, e := range errs {
		err := fmt.Sprintf("%v | %v",
			fmt.Sprintf("%s%s%s", RED, fmt.Sprintf("- SWW %s:%v", act, i+1), NOCOLOR),
			e.Error(),
		)

		text.Println(err)
	}
}

// Spinner generates static style notya spinner.
func Spinner() *spinner.Spinner {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Color("yellow")
	return s
}

// PrintServices logs given service names by provided color level.
func PrintServices(c string, services []string) {
	for _, s := range services {
		printable := fmt.Sprintf(" • %s", fmt.Sprintf("%s%s%s", c, s, NOCOLOR))
		text.Println(printable)
	}
}
