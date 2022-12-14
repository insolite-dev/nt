//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package models_test

import (
	"testing"

	"github.com/insolite-dev/notya/lib/models"
	"github.com/insolite-dev/notya/lib/services"
)

func TestToJSONofNote(t *testing.T) {
	tests := []struct {
		model    models.Note
		expected map[string]interface{}
	}{
		{
			model: models.Note{
				Title: "mock-title.txt",
				Path:  map[string]string{services.LOCAL.ToStr(): "~/mock-title.txt"},
				Body:  "empty",
			},
			expected: map[string]interface{}{
				"title": "mock-title.txt",
				"path":  map[string]string{services.LOCAL.ToStr(): "~/mock-title.txt"},
				"body":  "empty",
			},
		},
	}

	for _, td := range tests {
		got := td.model.ToJSON()
		for key, value := range td.expected {
			if got[key] != value {
				t.Errorf("NoteToJSON's sum was different: Want: %v | Got: %v", got, td.expected)
			}
		}
	}
}

func TestNoteToNode(t *testing.T) {
	tests := []struct {
		note     models.Note
		expected models.Node
	}{
		{
			note:     models.Note{},
			expected: models.Node{},
		},
		{
			note:     models.Note{Title: "title", Path: map[string]string{services.LOCAL.ToStr(): "~/title"}},
			expected: models.Node{Title: "title", Path: map[string]string{services.LOCAL.ToStr(): "~/title"}},
		},
	}

	for _, td := range tests {
		got := td.note.ToNode()
		path := got.GetPath(services.LOCAL.ToStr())
		if got.Title != td.expected.Title || path != td.expected.GetPath(services.LOCAL.ToStr()) {
			t.Errorf("Sum was different of [Note-to-Node] function: Want: %v | Got: %v", td.expected, got)
		}
	}
}
