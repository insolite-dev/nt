//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package commands

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/insolite-dev/notya/assets"
	"github.com/insolite-dev/notya/lib/models"
	"github.com/insolite-dev/notya/pkg"
	"github.com/spf13/cobra"
)

var cutCommand = &cobra.Command{
	Use:   "cut",
	Short: "Cut the file | copies the file and saves it data to clipboard",
	Run:   runCutCommand,
}

func initCutCommand() {
	appCommand.AddCommand(cutCommand)
}

// runCutCommand runs appropriate service commands to cut the note file.
func runCutCommand(cmd *cobra.Command, args []string) {
	determineService()

	if len(args) > 0 {
		cutAndFinish(models.Note{Title: args[0]})
		return
	}

	loading.Start()
	// Generate array of all node names.
	_, nodeNames, err := service.GetAll("", "file", models.NotyaIgnoreFiles)
	loading.Stop()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask for node selection.
	var selected string
	survey.AskOne(
		assets.ChooseNodePrompt("note", "cut", nodeNames),
		&selected,
	)

	cutAndFinish(models.Note{Title: selected})
}

func cutAndFinish(note models.Note) {
	if len(note.Title) == 0 {
		os.Exit(-1)
		return
	}

	loading.Start()
	if _, err := service.Cut(note); err != nil {
		loading.Stop()
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}
	loading.Stop()
}
