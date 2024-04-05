//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package commands

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/insolite-dev/nt/assets"
	"github.com/insolite-dev/nt/lib/models"
	"github.com/insolite-dev/nt/pkg"
	"github.com/spf13/cobra"
)

var whereCommand = &cobra.Command{
	Use:     "where",
	Aliases: []string{"path", "wh"},
	Short:   "View the path of file or folder",
	Run:     runWhereCommand,
}

func initWhereCommand() {
	appCommand.AddCommand(whereCommand)
}

func runWhereCommand(cmd *cobra.Command, args []string) {
	determineService()

	if len(args) > 0 {
		note, err := service.View(models.Note{Title: args[0]})
		loading.Stop()

		if err != nil {
			pkg.Alert(pkg.ErrorL, err.Error())
		} else {
			pkg.PrintPath((*note).ToNode())
		}

		return
	}

	nodes, noteNames, err := service.GetAll("", "", models.NotyaIgnoreFiles)
	loading.Stop()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask for note selection.
	var selected string
	survey.AskOne(
		assets.ChooseNodePrompt("note", "view path", noteNames),
		&selected,
	)

	for _, n := range nodes {
		if n.Title == selected {
			pkg.PrintPath(n)
		}
	}
}
