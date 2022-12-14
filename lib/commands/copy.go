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

// copyCommand is a command model which used to copy files.
var copyCommand = &cobra.Command{
	Use:     "copy",
	Aliases: []string{"c"},
	Short:   "Copy file's body to clipboard",
	Run:     runCopyCommand,
}

// initCopyCommand initializes copyCommand to the main application command.
func initCopyCommand() {
	appCommand.AddCommand(copyCommand)
}

// runCopyCommand runs appropriate service commands to copy note data to clipboard.
func runCopyCommand(cmd *cobra.Command, args []string) {
	determineService()

	if len(args) > 0 {
		copyAndFinish(models.Note{Title: args[0]})
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
		assets.ChooseNodePrompt("note", "copy", nodeNames),
		&selected,
	)

	copyAndFinish(models.Note{Title: selected})
}

func copyAndFinish(note models.Note) {
	if len(note.Title) == 0 {
		os.Exit(-1)
		return
	}

	loading.Start()
	if err := service.Copy(note); err != nil {
		loading.Stop()
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}
	loading.Stop()
}
