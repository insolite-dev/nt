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

// createCommand, is a command model which used to create new notes or files.
var createCommand = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "make"},
	Short:   "Create new note",
	Run:     runCreateCommand,
}

// initCreateCommand sets flags of command, and adds it to main application command.
func initCreateCommand() {
	appCommand.AddCommand(createCommand)
}

// runCreateCommand runs appropriate service commands to create new note.
func runCreateCommand(cmd *cobra.Command, args []string) {
	createAnswers := pkg.CreateAnswers{}

	// Start asking create command questions.
	if err := survey.Ask(
		pkg.CreateNoteQuestions,
		&createAnswers,
		survey.WithIcons(pkg.SurveyIconsConfig),
	); err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Create new note-file by [note].
	note, err := service.Create(models.Note{Title: createAnswers.Title})
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	if createAnswers.EditNote {
		// Open created note-file to edit it.
		if err := service.Open(*note); err != nil {
			pkg.Alert(pkg.ErrorL, err.Error())
			return
		}
	}

	pkg.Alert(pkg.SuccessL, "Created new note: "+note.Title)
}
