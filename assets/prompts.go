//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package assets

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/insolite-dev/notya/lib/models"
)

// ChooseNodePrompt is a prompt interface for tui file or folder choosing bar.
func ChooseNodePrompt(node, act string, options []string) *survey.Select {
	return &survey.Select{
		Message: fmt.Sprintf("Choose a %v to %v:", node, act),
		Options: options,
	}
}

// ChooseRemotePrompt is a prompt interface for tui remote service choosing bar.
func ChooseRemotePrompt(services []string) *survey.Select {
	return &survey.Select{
		Message: "Choose remote service:",
		Options: services,
	}
}

// CreatePromptQuestion is a question list for create command.
var CreatePromptQuestion = []*survey.Question{
	{
		Prompt: &survey.Input{
			Message: "Title",
			Help:    "Append to your note any title you want, and then complete file name with special file name type | e.g: new_note.md",
		},
		Validate: survey.MinLength(1),
	},
}

// Mkdir is a question list for mkdir command.
var MkdirPromptQuestion = []*survey.Question{
	{
		Prompt:   &survey.Input{Message: "Title"},
		Validate: survey.MinLength(1),
	},
}

// OpenViaEditorPromt is a confirm prompt for editor editing.
var OpenViaEditorPromt = &survey.Confirm{
	Message: "Wanna open with editor?",
	Help:    "Do you want to open note with your editor?",
	Default: false,
}

// NewNamePrompt is a input prompt for rename command.
func NewNamePrompt(d string) *survey.Input {
	return &survey.Input{Message: "New name: ", Default: d}
}

// SettingsEditPromptQuestions is a question list for settings' edit sub-command.
func SettingsEditPromptQuestions(defaultSettings models.Settings) []*survey.Question {
	return []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Default: defaultSettings.Name,
				Message: "App Name",
				Help:    "Customize your application env's name",
			},
			Validate: survey.MinLength(1),
		},
		{
			Name: "editor",
			Prompt: &survey.Input{
				Default: defaultSettings.Editor,
				Message: "Editor",
				Help:    "Editor for notya. --> vim/nvim/code/code-insiders ...",
			},
			Validate: survey.MinLength(1),
		},
		{
			Name: "local_path",
			Prompt: &survey.Input{
				Default: defaultSettings.LocalPath,
				Message: "Local Path",
				Help:    "Local path of notya base working directory",
			},
			Validate: survey.MinLength(1),
		},
		{
			Name: "fire_project_id",
			Prompt: &survey.Input{
				Default: defaultSettings.FirebaseProjectID,
				Message: "Firebase Project ID",
				Help:    "Project ID of your (integrated-with-notya) firebase project",
			},
		},
		{
			Name: "fire_account_key",
			Prompt: &survey.Input{
				Default: defaultSettings.FirebaseAccountKey,
				Message: "Firebase Key File",
				Help:    "Path of firebase service key file",
			},
		},
		{
			Name: "fire_collection",
			Prompt: &survey.Input{
				Default: defaultSettings.FirebaseCollection,
				Message: "Firebase Collection",
				Help:    "Main notya collection name in firestore",
			},
		},
	}
}

// MoveNotesPrompt is a confirm prompt for setting's move-note functionality.
var MoveNotesPrompt = &survey.Confirm{
	Message: "Move notes",
	Help:    "Do you wanna move old notes to new path?",
	Default: false,
}
