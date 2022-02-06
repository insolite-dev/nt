// Copyright 2021-2022 present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package assets_test

import (
	"errors"
	"testing"

	"github.com/anonistas/notya/assets"
)

func TestNotExists(t *testing.T) {
	tests := []struct {
		testname string
		path     string
		expected error
	}{
		{
			testname: "should generate not exists error",
			path:     "test/path.txt",
			expected: errors.New("File does not exists at: test/path.txt"),
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := assets.NotExists(td.path)
			if got.Error() != td.expected.Error() {
				t.Errorf("Sum of NotExists was different: Want: %v, Got: %v", td.expected, got)
			}
		})
	}
}

func TestAlreadyExists(t *testing.T) {
	tests := []struct {
		testname string
		path     string
		expected error
	}{
		{
			testname: "should generate already exists error",
			path:     "test/path.txt",
			expected: errors.New("A file already exists at: test/path.txt, please provide a unique title"),
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := assets.AlreadyExists(td.path)
			if got.Error() != td.expected.Error() {
				t.Errorf("Sum of AlreadyExists was different: Want: %v, Got: %v", td.expected, got)
			}
		})
	}
}
