// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package models_test

import (
	"testing"

	"github.com/anonistas/notya/lib/models"
)

func TestInitSettings(t *testing.T) {
	tests := []struct {
		testname string
		expected models.Settings
	}{
		{
			testname: "should return initial settings properly",
			expected: models.Settings{Editor: models.DefaultEditor, LocalPath: models.DefaultLocalPath},
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := models.InitSettings("notya")

			if got.Editor != td.expected.Editor || got.LocalPath != td.expected.LocalPath {
				t.Errorf("InitSettings's sum was different: Want: %v | Got: %v", got, td.expected)
			}
		})
	}
}

func TestToByte(t *testing.T) {
	tests := []struct {
		testname       string
		model          models.Settings
		expectedLength int
	}{
		{
			testname:       "should return initial settings properly",
			model:          models.Settings{Editor: "mvim"},
			expectedLength: 33,
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := td.model.ToByte()

			if len(got) != td.expectedLength {
				t.Errorf("ToByte's length sum was different: Want: %v | Got: %v", len(got), td.expectedLength)
			}
		})
	}
}

func TestFromJSON(t *testing.T) {
	tests := []struct {
		testname      string
		argumentValue string
		expected      models.Settings
	}{
		{
			testname:      "should return initial settings properly",
			argumentValue: `{"editor": "vi"}`,
			expected:      models.Settings{Editor: models.DefaultEditor},
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := models.FromJSON(td.argumentValue)

			if got.Editor != td.expected.Editor {
				t.Errorf("FromJSON's sum was different: Want: %v | Got: %v", got, td.expected)
			}
		})
	}
}
