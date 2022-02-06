// Copyright 2021-2022 present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package assets

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/anonistas/notya/lib/models"
)

// ChooseNotePrompt is a prompt interface for note choosing.
func ChooseNotePrompt(act string, options []string) *survey.Select {
	return &survey.Select{
		Message: fmt.Sprintf("Choose a note to %v:", act),
		Options: options,
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

// OpenViaEditorPromt is a confirm prompt for editor editing.
var OpenViaEditorPromt = &survey.Confirm{
	Message: "Wanna open with editor?",
	Help:    "Do you want to open note with your editor?",
	Default: false,
}

// NewNamePrompt is a input prompt for rename command.
func NewNamePrompt(d string) *survey.Input {
	return &survey.Input{
		Message: "New name: ",
		Help:    "Enter new note name/title (don't forget putting type of it, like: `renamed_note.txt`)",
		Default: d,
	}
}

// SettingsEditPromptQuestions is a question list for settings' edit sub-command.
func SettingsEditPromptQuestions(defaultSettings models.Settings) []*survey.Question {
	return []*survey.Question{
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
		},
	}
}

// MoveNotesPrompt is a confirm prompt for setting's move-note functionality.
var MoveNotesPrompt = &survey.Confirm{
	Message: "Move notes",
	Help:    "Do you wanna move old notes to new path?",
	Default: false,
}
