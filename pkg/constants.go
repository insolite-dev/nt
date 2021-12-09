package pkg

import "github.com/AlecAivazis/survey/v2"

// Version is current version of application.
const Version = "v1.0.0"

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

// CreateAnswers for the survey's answers for `create` command.
type CreateAnswers struct {
	Title    string
	EditNote bool `survey:"edit-note"`
}

var (
	// CreateNoteQuestions is a list of questions for create command.
	CreateNoteQuestions = []*survey.Question{
		{
			Name: "title",
			Prompt: &survey.Input{
				Message: "Enter name of new note: ",
				Help:    "Append to your note any name you want  and then, complete file name with special file name type | e.g: new_note.md",
			},
			Validate: survey.MinLength(1),
		},
		{
			Name: "edit-note",
			Prompt: &survey.Confirm{
				Message: "Do you wanna open note with Vi/Vim, to edit file?",
				Default: true,
			},
		},
	}
)
