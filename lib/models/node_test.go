// Copyright 2022-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package models_test

import (
	"testing"

	"github.com/anonistas/notya/lib/models"
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
			node:     models.Node{Title: "file.txt", Path: "~/file.txt"},
			expected: models.Note{Title: "file.txt", Path: "~/file.txt"},
		},
	}

	for _, td := range tests {
		got := td.node.ToNote()
		if got.Title != td.expected.Title || got.Path != td.expected.Path {
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
			node:     models.Node{Title: "folder/", Path: "~/folder/"},
			expected: models.Folder{Title: "folder/", Path: "~/folder/"},
		},
	}

	for _, td := range tests {
		got := td.node.ToFolder()
		if got.Title != td.expected.Title || got.Path != td.expected.Path {
			t.Errorf("Sum was different of [Node-To-Folder] function: Want: %v | Got: %v", td.expected, got)
		}
	}
}

func TestStructAsFolder(t *testing.T) {
	tests := []struct {
		node     models.Node
		expected models.Node
	}{
		{node: models.Node{}, expected: models.Node{}},
		{
			node:     models.Node{Title: "folder", Path: "~/folder"},
			expected: models.Node{Title: "folder/", Path: "~/folder/"},
		},
	}

	for _, td := range tests {
		got := td.node.StructAsFolder()
		if got.Title != td.expected.Title || got.Path != td.expected.Path {
			t.Errorf("Sum was different of [StructAsFolder] function: Want: %v | Got: %v", td.expected, got)
		}
	}
}

func TestStructAsNote(t *testing.T) {
	tests := []struct {
		node     models.Node
		expected models.Node
	}{
		{node: models.Node{}, expected: models.Node{}},
		{
			node:     models.Node{Title: "note.txt/", Path: "~/note.txt/"},
			expected: models.Node{Title: "note.txt", Path: "~/note.txt"},
		},
	}

	for _, td := range tests {
		got := td.node.StructAsNote()
		if got.Title != td.expected.Title || got.Path != td.expected.Path {
			t.Errorf("Sum was different of [StructAsNote] function: Want: %v | Got: %v", td.expected, got)
		}
	}
}
