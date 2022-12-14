//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package assets_test

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/insolite-dev/notya/assets"
)

func TestChoseNotePrompt(t *testing.T) {
	type arguments struct {
		node, msg string
		options   []string
	}

	tests := []struct {
		testname string
		args     arguments
		expected survey.Select
	}{
		{
			testname: "should generate choosing-note prompt properly",
			args: arguments{
				node:    "note",
				msg:     "edit",
				options: []string{"1", "2", "3"},
			},
			expected: survey.Select{
				Message: "Choose a note to edit:",
				Options: []string{"1", "2", "3"},
			},
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := assets.ChooseNodePrompt(td.args.node, td.args.msg, td.args.options)

			// Closure function to check if options are different or not.
			var isDiffArr = func() bool {
				var a1, a2 = got.Options, td.expected.Options
				if len(a1) != len(a2) {
					return true
				}

				for i := 0; i < len(a1); i++ {
					if a1[i] != a2[i] {
						return true
					}
				}

				return false
			}()

			if got.Message != td.expected.Message || isDiffArr {
				t.Errorf("Sum of ChooseNotePrompt was different: Want: %v | Got: %v", td.expected, got)
			}
		})
	}
}

func TestChooseRemotePrompt(t *testing.T) {
	tests := []struct {
		services []string
		expected *survey.Select
	}{
		{
			services: []string{"LOCAL", "FIREBASE"},
			expected: &survey.Select{
				Message: "Choose remote service:",
				Options: []string{"LOCAL", "FIREBASE"},
			},
		},
	}

	for _, td := range tests {
		got := assets.ChooseRemotePrompt(td.services)

		// Closure function to check if options are different or not.
		var isDiffArr = func() bool {
			var a1, a2 = got.Options, td.expected.Options
			if len(a1) != len(a2) {
				return true
			}

			for i := 0; i < len(a1); i++ {
				if a1[i] != a2[i] {
					return true
				}
			}

			return false
		}()

		if isDiffArr || got.Message != td.expected.Message {
			t.Errorf("Sum of ChooseRemotePrompt was different: Want: %v | Got: %v", td.expected, got)
		}
	}
}

func TestNewNamePrompt(t *testing.T) {
	tests := []struct {
		testname     string
		defaultValue string
		expected     survey.Input
	}{
		{
			testname:     "should generate new-name-prompt properly",
			defaultValue: "default-name",
			expected: survey.Input{
				Message: "New name: ",
				Default: "default-name",
			},
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := assets.NewNamePrompt(td.defaultValue)

			if got.Message != td.expected.Message || got.Help != td.expected.Help || got.Default != td.expected.Default {
				t.Errorf("Sum of NewNamePrompt was different: Want: %v | Got: %v", td.expected, got)
			}
		})
	}
}
