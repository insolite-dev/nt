//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package assets

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
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

// MoveNotesPrompt is a confirm prompt for setting's move-note functionality.
var MoveNotesPrompt = &survey.Confirm{
	Message: "Move notes",
	Help:    "Do you wanna move old notes to new path?",
	Default: false,
}

// FirebaseRemoteConnectPromptQuestion is a question list that fills up
// required values for firebase remote connection.
// Used in Remote command's connect subcommand.
var FirebaseRemoteConnectPromptQuestion = []*survey.Question{
	{
		Name: "fire_project_id",
		Prompt: &survey.Input{
			Message: "Firebase Project ID",
			Help:    "The project ID of your Firebase project.",
		},
		Validate: survey.MinLength(1),
	},
	{
		Name: "fire_account_key",
		Prompt: &survey.Input{
			Message: "Firebase Account Key",
			Help:    "The Firebase Admin SDK private key file path. Must be given a full path, like: /Users/john-doe/notya/account_key.json.",
		},
		Validate: survey.MinLength(5),
	},
	{
		Name: "fire_collection",
		Prompt: &survey.Input{
			Message: "Firebase Collection",
			Help:    "A name of collection for notes, from your firebase project's firestore.",
		},
		Validate: survey.MinLength(1),
	},
}
