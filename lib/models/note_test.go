// Copyright 2022-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package models_test

import (
	"testing"

	"github.com/anonistas/notya/lib/models"
)

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
			note:     models.Note{Title: "title", Path: "~/title"},
			expected: models.Node{Title: "title", Path: "~/title"},
		},
	}

	for _, td := range tests {
		got := td.note.ToNode()
		if got.Title != td.expected.Title || got.Path != td.expected.Path {
			t.Errorf("Sum was different of [Note-to-Node] function: Want: %v | Got: %v", td.expected, got)
		}
	}
}
