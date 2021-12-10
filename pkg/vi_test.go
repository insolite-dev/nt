// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package pkg_test

import (
	"errors"
	"testing"

	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
)

func TestOpenFileWithVI(t *testing.T) {
	type utilArgs struct {
		filename       string
		stdargs        models.StdArgs
		deleteFileFunc func(filename string)
		createFileFunc func(filename string)
	}

	tests := []struct {
		testName string
		ua       utilArgs
		expected error
	}{
		{
			testName: "should open created exiting file properly",
			ua: utilArgs{
				filename: "test_file.txt",
				stdargs:  models.StdArgs{},
				deleteFileFunc: func(filename string) {
					pkg.Delete(filename)
				},
				createFileFunc: func(filename string) {
					pkg.NewFile(filename, []byte{})
				},
			},

			expected: errors.New("exit status 1"),
		},
	}

	for _, td := range tests {
		t.Run(td.testName, func(t *testing.T) {
			td.ua.createFileFunc(td.ua.filename)

			got := pkg.OpenFileWithVI(td.ua.filename, td.ua.stdargs)
			if got != td.expected && got.Error() != td.expected.Error() {
				t.Errorf("Sum was different, Got: %v | Want: %v", got, td.expected)
			}

			td.ua.deleteFileFunc(td.ua.filename)
		})
	}
}
