// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package pkg_test

import (
	"errors"
	"os"
	"testing"

	"github.com/anonistas/notya/pkg"
)

func TestNotyaPWD(t *testing.T) {
	// Take current working directory first.
	currentHomeDir, _ := os.UserHomeDir()

	type expected struct {
		res string
		err error
	}

	tests := []struct {
		testName string
		exp      expected
	}{
		{
			testName: "should get right notya notes path",
			exp:      expected{currentHomeDir + "/notya-notes/", nil},
		},
	}

	for _, td := range tests {
		gotRes, gotErr := pkg.NotyaPWD()
		if gotErr != td.exp.err {
			t.Errorf("Path err sum was different: Got: %v | Want: %v", gotErr, td.exp.err)
		}

		if *gotRes != td.exp.res {
			t.Errorf("Path res sum was different: Got: %v | Want: %v", gotRes, td.exp.res)
		}
	}
}

func TestFileExists(t *testing.T) {
	type closures struct {
		creating func(name string)
		deleting func(name string)
	}

	tests := []struct {
		testName string
		filename string
		c        closures
		expected bool
	}{
		{
			"should check file not exists, properly",
			"test.txt",
			closures{creating: func(name string) {}, deleting: func(name string) {}},
			false,
		},
		{
			"should check file exists, properly",
			"test.txt",
			closures{
				creating: func(name string) { pkg.NewFile(name, []byte{}) },
				deleting: func(name string) { pkg.Delete(name) },
			},
			true,
		},
	}

	for _, td := range tests {
		td.c.creating(td.filename)

		got := pkg.FileExists(td.filename)
		if got != td.expected {
			t.Errorf("FileExists sum was different: Got: %v | Want: %v", got, td.expected)
		}

		td.c.deleting(td.filename)
	}
}

func TestNewFile(t *testing.T) {
	type args struct {
		filename string
		filebody []byte
	}

	tests := []struct {
		testName string
		a        args
		expected error
	}{
		{
			"should create new file properly",
			args{"test.txt", []byte{}},
			nil,
		},
	}

	for _, td := range tests {
		got := pkg.NewFile(td.a.filename, td.a.filebody)

		defer pkg.Delete(td.a.filename)

		if got != td.expected {
			t.Errorf("NewFile sum was different: Got: %v | Want: %v", got, td.expected)
		}

	}
}

func TestNewFolder(t *testing.T) {
	tests := []struct {
		testName      string
		foldernameArg string
		deleteFunc    func(foldername string)
		expected      error
	}{
		{
			testName:      "should create new folder properly | without deleting it",
			foldernameArg: "test_folder",
			deleteFunc:    func(foldername string) {},
			expected:      nil,
		},
		{
			testName:      "should alert error on trying to create already created folder",
			foldernameArg: "test_folder",
			deleteFunc: func(foldername string) {
				pkg.Delete(foldername)
			},
			expected: errors.New("mkdir test_folder: file exists"),
		},
	}

	for _, td := range tests {
		got := pkg.NewFolder(td.foldernameArg)

		defer td.deleteFunc(td.foldernameArg)

		if got != td.expected && got.Error() != td.expected.Error() {
			t.Errorf("NewFolder sum was different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		fileName       string
		createFileFunc func(filename string)
	}

	tests := []struct {
		testName string
		a        args
		expected interface{}
	}{
		{
			testName: "should delete exiting folder properly",
			a: args{"test_folder", func(filename string) {
				pkg.NewFolder(filename)
			}},
			expected: nil,
		},
		{
			testName: "should alert error, on trying deleting non-exiting file",
			a:        args{"test_folder", func(filename string) {}},
			expected: "remove test_folder: no such file or directory",
		},
	}

	for _, td := range tests {
		td.a.createFileFunc(td.a.fileName)

		got := pkg.Delete(td.a.fileName)
		if got != td.expected && got.Error() != td.expected {
			t.Errorf("NewFolder sum was different: Got: %v | Want: %v", got, td.expected)
		}

	}

}
