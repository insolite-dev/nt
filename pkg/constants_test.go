//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package pkg_test

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/insolite-dev/notya/pkg"
)

func TestSurveyIconsConfig(t *testing.T) {
	type expected struct {
		questionFormat, questionText string
		helpFormat, helpText         string
		errorFormat, errorText       string
	}

	tests := []struct {
		testName string
		e        expected
	}{
		{
			testName: "config should be expected properly",
			e: expected{
				questionFormat: "cyan",
				questionText:   "[?]",
				helpFormat:     "blue",
				helpText:       "Help ->",
				errorFormat:    "yellow",
				errorText:      "Warning ->",
			},
		},
	}

	for _, td := range tests {
		t.Run(td.testName, func(t *testing.T) {
			var got survey.IconSet
			var generateGot = func(setIcon func(*survey.IconSet)) {
				setIcon(&got)
			}

			generateGot(pkg.SurveyIconsConfig)

			if got.Question.Format != td.e.questionFormat || got.Question.Text != td.e.questionText {
				t.Errorf("Sum of question is different: Got: %v | Want: %v", got.Question, survey.Icon{Format: td.e.questionFormat, Text: td.e.questionText})
			}

			if got.Help.Format != td.e.helpFormat || got.Help.Text != td.e.helpText {
				t.Errorf("Sum of help is different: Got: %v | Want: %v", got.Help, survey.Icon{Format: td.e.helpFormat, Text: td.e.helpText})
			}

			if got.Error.Format != td.e.errorFormat || got.Error.Text != td.e.errorText {
				t.Errorf("Sum of error is different: Got: %v | Want: %v", got.Error, survey.Icon{Format: td.e.errorFormat, Text: td.e.errorText})
			}
		})
	}
}
