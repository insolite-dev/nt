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

// viewCommand, is a command model which used to view body of notes or files.
var viewCommand = &cobra.Command{
	Use:     "view",
	Aliases: []string{"show", "read"},
	Short:   "View full note data",
	Run:     runViewCommand,
}

// initViewCommand adds viewCommand to main application command.
func initViewCommand() {
	appCommand.AddCommand(viewCommand)
}

// runViewCommand runs appropriate service commands to log full note data.
func runViewCommand(cmd *cobra.Command, args []string) {
	// Take note title from arguments. If it's provided.
	if len(args) > 0 {
		note, err := service.ViewNote(models.Note{Title: args[0]})
		if err != nil {
			pkg.Alert(pkg.ErrorL, err.Error())
			return
		}

		pkg.ShowNote(*note)
		return
	}

	// If note name wasn't provided by arguments,
	// Shows select list by all notes list.
	// See: https://github.com/AlecAivazis/survey#select

	// Generate array of all notes' names.
	notes, err := service.GetAll()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask for note selection.
	selected := ""
	prompt := &survey.Select{
		Message: "Choose a note to view:",
		Options: notes,
	}
	survey.AskOne(prompt, &selected)

	// Get selected note.
	note, viewErr := service.ViewNote(models.Note{Title: selected})
	if viewErr != nil {
		pkg.Alert(pkg.ErrorL, viewErr.Error())
		return
	}

	pkg.ShowNote(*note)
}
