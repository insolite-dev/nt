package pkg

import "github.com/AlecAivazis/survey/v2"

var (
	// Custom configuration for survey icons and colors.
	// See: https://github.com/mgutz/ansi#style-format
	SurveyIconsConfig = func(icons *survey.IconSet) {
		icons.Question.Format = "cyan"
		icons.Question.Text = "[?]"
		icons.Help.Format = "blue"
		icons.Help.Text = "Help ->"
		icons.Error.Format = "yellow"
		icons.Error.Text = "Note ->"
	}
)

var (
	// OpenNoteToEdit is a early created survey question
	// Basically used on after creating new note.
	OpenNoteToEdit = survey.Question{
		Name: "edit-note",
		Prompt: &survey.Confirm{
			Message: "Do you wanna open note with Vi/Vim, to edit file?",
			Default: true,
		},
	}
)
