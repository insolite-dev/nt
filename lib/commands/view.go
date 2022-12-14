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

// viewCommand is a command model which used to view metadata of note.
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
	determineService()

	loading.Start()

	// Take note title from arguments. If it's provided.
	if len(args) > 0 {
		note, err := service.View(models.Note{Title: args[0]})
		loading.Stop()

		if err != nil {
			pkg.Alert(pkg.ErrorL, err.Error())
		} else {
			pkg.PrintNote(*note)
		}

		return
	}

	// Generate array of all note names.
	nodes, noteNames, err := service.GetAll("", "file", models.NotyaIgnoreFiles)
	loading.Stop()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask for note selection.
	var selected string
	survey.AskOne(
		assets.ChooseNodePrompt("note", "view", noteNames),
		&selected,
	)

	for _, n := range nodes {
		if n.Title == selected {
			pkg.PrintNote(n.ToNote())
		}
	}
}
