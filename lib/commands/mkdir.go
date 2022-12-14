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

// mkdirCommand is a command model that used to create new folders.
var mkdirCommand = &cobra.Command{
	Use:     "mkdir",
	Aliases: []string{"md"},
	Short:   "Create new working directory(folder)",
	Run:     runMkdirCommand,
}

// initMkdirCommand adds it to the main application command.
func initMkdirCommand() {
	appCommand.AddCommand(mkdirCommand)
}

// runMkdirCommand() runs appropriate service commands to create new folder.
func runMkdirCommand(cmd *cobra.Command, args []string) {
	determineService()

	var title string

	if len(args) > 0 { // Take folder's title from arguments, if it's provided.
		title = args[0]
	} else { // Ask for the title of folder.
		survey.Ask(assets.MkdirPromptQuestion, &title)
	}

	loading.Start()

	if len(title) == 0 {
		os.Exit(-1)
		return
	}

	// Create new directory by given title.
	_, err := service.Mkdir(models.Folder{Title: title})

	loading.Stop()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}
}
