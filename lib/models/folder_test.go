// Copyright 2022-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package models_test

import (
	"testing"

	"github.com/anonistas/notya/lib/models"
)

func TestFolderToNode(t *testing.T) {
	tests := []struct {
		dir      models.Folder
		expected models.Node
	}{
		{
			dir:      models.Folder{},
			expected: models.Node{},
		},
		{
			dir:      models.Folder{Title: "folder/", Path: "~/folder"},
			expected: models.Node{Title: "folder/", Path: "~/folder"},
		},
	}

	for _, td := range tests {
		got := td.dir.ToNode()
		if got.Title != td.expected.Title || got.Path != td.expected.Path {
			t.Errorf("Sum was different of [Folder-To-Node] function: Want: %v | Got: %v", td.expected, got)
		}
	}
}
