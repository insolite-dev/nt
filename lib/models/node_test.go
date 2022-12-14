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

func TestToNote(t *testing.T) {
	tests := []struct {
		node     models.Node
		expected models.Note
	}{
		{
			node:     models.Node{},
			expected: models.Note{},
		},
		{
			node:     models.Node{Title: "file.txt", Path: map[string]string{services.LOCAL.ToStr(): "~/file.txt"}},
			expected: models.Note{Title: "file.txt", Path: map[string]string{services.LOCAL.ToStr(): "~/file.txt"}},
		},
	}

	for _, td := range tests {
		got := td.node.ToNote()
		path := got.GetPath(services.LOCAL.ToStr())
		if got.Title != td.expected.Title || path != td.expected.GetPath(path) {
			t.Errorf("Sum was different of [Node-To-Note] function: Want: %v | Got: %v", td.expected, got)
		}
	}
}

func TestToFolder(t *testing.T) {
	tests := []struct {
		node     models.Node
		expected models.Folder
	}{
		{
			node:     models.Node{},
			expected: models.Folder{},
		},
		{
			node:     models.Node{Title: "folder/", Path: map[string]string{services.LOCAL.ToStr(): "~/folder/"}},
			expected: models.Folder{Title: "folder/", Path: map[string]string{services.LOCAL.ToStr(): "~/folder/"}},
		},
	}

	for _, td := range tests {
		got := td.node.ToFolder()
		path := got.GetPath(services.LOCAL.ToStr())
		if got.Title != td.expected.Title || path != td.expected.GetPath(services.LOCAL.ToStr()) {
			t.Errorf("Sum was different of [Node-To-Folder] function: Want: %v | Got: %v", td.expected, got)
		}
	}
}
