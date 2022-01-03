// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
	"github.com/spf13/cobra"
)

// createCommand is a command model that used to create new notes or files.
var createCommand = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "make"},
	Short:   "Create new note",
	Run:     runCreateCommand,
}

// initCreateCommand adds it to the main application command.
func initCreateCommand() {
	appCommand.AddCommand(createCommand)
}

// runCreateCommand runs appropriate service commands to create new note.
func runCreateCommand(cmd *cobra.Command, args []string) {
	// Take new note's title from arguments, if it's provided.
	if len(args) > 0 {
		title := args[0]
		createAndFinish(title)
		return
	}

	// Ask for title of new note.
	var title string
	question := []*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Enter name of new note: ",
				Help:    "Append to your note any name you want and then, complete file name with special file name type | e.g: new_note.md",
			},
			Validate: survey.MinLength(1),
		},
	}
	survey.Ask(question, &title)

	createAndFinish(title)
}

// createAndFinish asks to edit note and finishes creating loop.
func createAndFinish(title string) {
	// Create new note-file by given title.
	note, err := service.Create(models.Note{Title: title})
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask for, open or not created note with editor.
	var openNote bool
	prompt := &survey.Confirm{
		Message: "Do you wanna open note with your editor?",
		Default: false,
	}
	survey.AskOne(prompt, &openNote)

	if openNote {
		// Open created note-file to edit it.
		if err := service.Open(*note); err != nil {
			pkg.Alert(pkg.ErrorL, err.Error())
			return
		}
	}

	pkg.Alert(pkg.SuccessL, "Created new note: "+note.Title)
}
