// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package pkg_test

import (
	"errors"
	"os"
	"testing"

	"github.com/anonistas/notya/lib/models"
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
			exp:      expected{currentHomeDir + "/notya", nil},
		},
	}

	for _, td := range tests {
		gotRes, gotErr := pkg.NotyaPWD(models.Settings{LocalPath: "notya"})
		if gotErr != td.exp.err {
			t.Errorf("Path err sum was different: Got: %v | Want: %v", gotErr, td.exp.err)
		}

		if *gotRes != td.exp.res {
			t.Errorf("Path res sum was different: Got: %v | Want: %v", *gotRes, td.exp.res)
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
				creating: func(name string) { pkg.WriteNote(name, []byte{}) },
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

func TestWriteNote(t *testing.T) {
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
		got := pkg.WriteNote(td.a.filename, td.a.filebody)

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

func TestReadBody(t *testing.T) {
	type expected struct {
		err error
	}

	test := []struct {
		testName     string
		fileName     string
		creatingFunc func(filename string)
		deletingFunc func(filename string)
		e            expected
	}{
		{
			testName: "should read file properly",
			fileName: "test_file.txt",
			creatingFunc: func(filename string) {
				pkg.WriteNote(filename, []byte{})
			},
			deletingFunc: func(filename string) {
				pkg.Delete(filename)
			},
			e: expected{err: nil},
		},
	}

	for _, td := range test {
		td.creatingFunc(td.fileName)

		t.Run(td.testName, func(t *testing.T) {
			_, err := pkg.ReadBody(td.fileName)
			if err != td.e.err {
				t.Errorf("ReadBody err sum was different, Got: %v | Want: %v", err, td.e.err)
			}
		})

		td.deletingFunc(td.fileName)
	}
}

func TestListDir(t *testing.T) {
	type expected struct {
		res []string
		err error
	}

	tests := []struct {
		testName     string
		folderName   string
		creatingFunc func(foldername string)
		deletingFunc func(foldername string)
		e            expected
	}{
		{
			testName:   "should list directory files properly",
			folderName: "test_folder",
			creatingFunc: func(foldername string) {
				pkg.NewFolder(foldername)
				pkg.WriteNote(foldername+"/test_file.txt", []byte{})
				pkg.WriteNote(foldername+"/test_file_1.txt", []byte{})
				pkg.WriteNote(foldername+"/expectable.txt", []byte{})
			},
			deletingFunc: func(foldername string) {
				pkg.Delete(foldername + "/test_file.txt")
				pkg.Delete(foldername + "/test_file_1.txt")
				pkg.Delete(foldername + "/expectable.txt")
				pkg.Delete(foldername)
			},
			e: expected{
				res: []string{"test_file.txt", "test_file_1.txt"},
				err: nil,
			},
		},
	}

	for _, td := range tests {
		td.creatingFunc(td.folderName)

		t.Run(td.testName, func(t *testing.T) {
			got, err := pkg.ListDir(td.folderName, "expectable.txt")
			if err != td.e.err {
				t.Errorf("ListDir's error sum was different, Got: %v | Want: %v", err, td.e.err)
			}

			if len(got) != len(td.e.res) {
				t.Errorf("ListDir's res length sum was different, Got: %v | Want: %v", len(got), len(td.e.res))
			}

			// Check each element of got
			for i := 0; i < len(got); i++ {
				if got[i] != td.e.res[i] {
					t.Errorf("ListDir's res [%v] item sum was different, Got: %v | Want: %v", i, err, td.e.err)
				}
			}
		})

		td.deletingFunc(td.folderName)
	}
}

func TestOpenViaEditor(t *testing.T) {
	type utilArgs struct {
		filename       string
		stdargs        models.StdArgs
		settings       models.Settings
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
				settings: models.InitSettings("notya"),
				deleteFileFunc: func(filename string) {
					pkg.Delete(filename)
				},
				createFileFunc: func(filename string) {
					pkg.WriteNote(filename, []byte{})
				},
			},
			expected: errors.New("signal: abort trap"),
		},
	}

	for _, td := range tests {
		t.Run(td.testName, func(t *testing.T) {
			td.ua.createFileFunc(td.ua.filename)

			got := pkg.OpenViaEditor(td.ua.filename, td.ua.stdargs, td.ua.settings)
			if got != td.expected && got.Error() != td.expected.Error() {
				t.Errorf("Sum was different, Got: %v | Want: %v", got, td.expected)
			}

			td.ua.deleteFileFunc(td.ua.filename)
		})
	}
}
