//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package models_test

import (
	"testing"

	"github.com/insolite-dev/nt/lib/models"
)

func TestInitSettings(t *testing.T) {
	tests := []struct {
		testname string
		expected models.Settings
	}{
		{
			testname: "should return initial settings properly",
			expected: models.Settings{Editor: models.DefaultEditor, NotesPath: models.DefaultLocalPath},
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := models.InitSettings("nt")

			if got.Editor != td.expected.Editor || got.NotesPath != td.expected.NotesPath {
				t.Errorf("InitSettings's sum was different: Want: %v | Got: %v", got, td.expected)
			}
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		testname       string
		model          models.Settings
		expectedLength int
	}{
		{
			testname:       "should return initial settings properly",
			model:          models.Settings{Editor: "mvim"},
			expectedLength: 56,
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := td.model.ToString()

			if len(got) != td.expectedLength {
				t.Errorf("ToString's length sum was different: Want: %v | Got: %v", td.expectedLength, len(got))
			}
		})
	}
}

func TestToJSON(t *testing.T) {
	tests := []struct {
		model    models.Settings
		expected map[string]interface{}
	}{
		{
			model: models.Settings{
				Name:               models.DefaultAppName,
				Editor:             models.DefaultEditor,
				NotesPath:          "~nt",
				FirebaseProjectID:  "nt",
				FirebaseAccountKey: "~nt/key.json",
				FirebaseCollection: "nt-notes",
			},
			expected: map[string]interface{}{
				"name":             models.DefaultAppName,
				"editor":           models.DefaultEditor,
				"notes_path":       "~nt",
				"fire_project_id":  "nt",
				"fire_account_key": "~nt/key.json",
				"fire_collection":  "nt-notes",
			},
		},
	}

	for _, td := range tests {
		got := td.model.ToJSON()

		for key, value := range td.expected {

			if got[key] != value {
				t.Errorf("SettingsToJSON's sum was different: Want: %v | Got: %v", got, td.expected)
			}
		}
	}
}

func TestDecodeSettings(t *testing.T) {
	tests := []struct {
		testname      string
		argumentValue string
		expected      models.Settings
	}{
		{
			testname:      "should generate settings model from json properly",
			argumentValue: `{"editor": "vi"}`,
			expected:      models.Settings{Editor: models.DefaultEditor},
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := models.DecodeSettings(td.argumentValue)

			if got.Editor != td.expected.Editor {
				t.Errorf("DecodeSettings's sum was different: Want: %v | Got: %v", got, td.expected)
			}
		})
	}
}

func TestFirePath(t *testing.T) {
	tests := []struct {
		model    models.Settings
		expected string
	}{
		{
			model:    models.Settings{},
			expected: "nt",
		},
		{
			model:    models.Settings{Name: "nt"},
			expected: "nt",
		},
		{
			model:    models.Settings{FirebaseCollection: "nt-notes", Name: "nt"},
			expected: "nt-notes",
		},
	}

	for _, td := range tests {
		got := td.model.FirePath()

		if got != td.expected {
			t.Errorf("FirePath's sum was different: Want: %v | Got: %v", got, td.expected)
		}
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		testname string
		settings models.Settings
		expected bool
	}{
		{
			testname: "should check settings validness correctly | [valid]",
			settings: models.InitSettings("/usr/mock/NotesPath"),
			expected: true,
		},
		{
			testname: "should check settings validness correctly | [invalid]",
			settings: models.Settings{},
			expected: false,
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := td.settings.IsValid()

			if got != td.expected {
				t.Errorf("IsValid sum was different: Want: %v | Got: %v", got, td.expected)
			}
		})
	}
}
