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

// createCommand is a command model that used to create new notes or files.
var createCommand = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create new node(file/folder)",
	Run:     runCreateCommand,
}

// providedFolderName is the value of folder flag.
var providedFolderName string

// initCreateCommand adds it to the main application command.
func initCreateCommand() {
	createCommand.Flags().StringVarP(
		&providedFolderName, "folder", "d", "",
		"Make a directory via create command",
	)

	appCommand.AddCommand(createCommand)
}

// runCreateCommand runs appropriate service commands to create new note.
func runCreateCommand(cmd *cobra.Command, args []string) {
	determineService()

	// Move direction to mkdir command.
	if providedFolderName != "" {
		runMkdirCommand(cmd, []string{providedFolderName})
		return
	}

	// Take new note's title from arguments, if it's provided.
	if len(args) > 0 {
		title := args[0]

		// Check for directory rep-slash at end of the title.
		// If it's provided, create command should switch  functionality
		// to mkdir command.
		if string(title[len(title)-1]) == "/" {
			runMkdirCommand(cmd, []string{title})
			return
		}

		createAndFinish(title)
		return
	}

	// Ask for title of new note.
	var title string
	survey.Ask(assets.CreatePromptQuestion, &title)

	createAndFinish(title)
}

// createAndFinish asks to edit note and finishes creating loop.
func createAndFinish(title string) {
	if len(title) == 0 {
		os.Exit(-1)
		return
	}

	loading.Start()
	note, err := service.Create(models.Note{Title: title})
	loading.Stop()

	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
		return
	}

	// Ask for, open or not created note with editor.
	var openNote bool
	survey.AskOne(assets.OpenViaEditorPromt, &openNote)

	if openNote {
		// Open created note-file to edit it.
		if err := service.Open(note.ToNode()); err != nil {
			pkg.Alert(pkg.ErrorL, err.Error())
			return
		}
	}
}
