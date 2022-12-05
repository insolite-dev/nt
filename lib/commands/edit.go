//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package commands

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/insolite-dev/notya/assets"
	"github.com/insolite-dev/notya/lib/models"
	"github.com/insolite-dev/notya/pkg"
	"github.com/spf13/cobra"
)

// editCommand is a command model which used to overwrite body of notes or files.
var editCommand = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"overwrite", "update"},
	Short:   "Edit/Update note data",
	Run:     runEditCommand,
}

// initEditCommand adds editCommand to main application command.
func initEditCommand() {
	appCommand.AddCommand(editCommand)
}

// runEditCommand runs appropriate service commands to edit/overwrite note data.
func runEditCommand(cmd *cobra.Command, args []string) {
	determineService()

	// Take note title from arguments. If it's provided.
	if len(args) > 0 {
		note := models.Note{Title: args[0]}

		if err := service.Open(note.ToNode()); err != nil {
			pkg.Alert(pkg.ErrorL, err.Error())
		}

		return
	}

	// Generate all node names.
	loading.Start()
	_, nodeNames, err := service.GetAll("", "file", models.NotyaIgnoreFiles)
	loading.Stop()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask for note selection.
	var selected string
	survey.AskOne(
		assets.ChooseNodePrompt("note", "edit", nodeNames),
		&selected,
	)

	// Open selected note-file.
	if err := service.Open(models.Node{Title: selected}); err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}
}
