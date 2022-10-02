//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package assets_test

import (
	"errors"
	"testing"

	"github.com/insolite-dev/notya/assets"
)

func TestNotExists(t *testing.T) {
	tests := []struct {
		testname   string
		path, node string
		expected   error
	}{
		{
			testname: "should generate not exists error without path",
			path:     "",
			node:     "File",
			expected: errors.New("File does not exists"),
		},
		{
			testname: "should generate not exists error with path",
			path:     "test/path/",
			node:     "Directory",
			expected: errors.New("Directory does not exists at: test/path/"),
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := assets.NotExists(td.path, td.node)
			if got.Error() != td.expected.Error() {
				t.Errorf("Sum of NotExists was different: Want: %v, Got: %v", td.expected, got)
			}
		})
	}
}

func TestAlreadyExists(t *testing.T) {
	tests := []struct {
		testname   string
		path, node string
		expected   error
	}{
		{
			testname: "should generate already exists error for note",
			node:     "file",
			path:     "test/path.txt",
			expected: errors.New("A file already exists at: test/path.txt, please provide a unique title"),
		},
		{
			testname: "should generate already exists error for folder",
			node:     "directory",
			path:     "test/new-folder/",
			expected: errors.New("A directory already exists at: test/new-folder/, please provide a unique title"),
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := assets.AlreadyExists(td.path, td.node)
			if got.Error() != td.expected.Error() {
				t.Errorf("Sum of AlreadyExists was different: Want: %v, Got: %v", td.expected, got)
			}
		})
	}
}

func TestCannotDoSth(t *testing.T) {
	tests := []struct {
		act, doc      string
		err, expected error
	}{
		{
			act: "fetch", doc: "note",
			err:      errors.New("sww"),
			expected: errors.New("Cannot fetch note | sww"),
		},
	}

	for _, td := range tests {
		got := assets.CannotDoSth(td.act, td.doc, td.err)
		if got.Error() != td.expected.Error() {
			t.Errorf("Sum of CannotDoSth was different: Want: %v, Got: %v", td.expected, got)
		}
	}

}
