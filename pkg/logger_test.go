//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package pkg_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/insolite-dev/notya/lib/models"
	"github.com/insolite-dev/notya/lib/services"
	"github.com/insolite-dev/notya/pkg"
)

func TestAlert(t *testing.T) {
	tests := []struct {
		testName string
		level    pkg.Level
		message  string
	}{
		{
			"should alert error-level message",
			pkg.ErrorL,
			"Cannot be created new note",
		},
		{
			"should alert success-level message",
			pkg.SuccessL,
			"New note created successfully",
		},
		{
			"should alert info-level message",
			pkg.InfoL,
			"Merging local files with db files took 3m",
		},
	}
	for _, td := range tests {
		t.Run(
			td.testName,
			func(t *testing.T) {
				pkg.Alert(td.level, td.message)
			},
		)
	}
}

func TestOutputLevel(t *testing.T) {
	tests := []struct {
		testName string
		level    pkg.Level
		expected string
	}{
		{
			"should send normal message",
			pkg.Level("nocolor-default"),
			fmt.Sprintf("%s%s%s", pkg.NOCOLOR, "", pkg.NOCOLOR),
		},
		{
			"should send success message",
			pkg.SuccessL,
			fmt.Sprintf("%s%s%s", pkg.GREEN, pkg.SUCCESS, pkg.NOCOLOR),
		},
		{
			"should send error message",
			pkg.ErrorL,
			fmt.Sprintf("%s%s%s", pkg.RED, pkg.ERROR, pkg.NOCOLOR),
		},
		{
			"should send info message",
			pkg.InfoL,
			fmt.Sprintf("%s%s%s", pkg.YELLOW, pkg.INFO, pkg.NOCOLOR),
		},
	}

	for _, td := range tests {
		t.Run(
			td.testName,
			func(t *testing.T) {
				got := pkg.OutputLevel(td.level)
				if got != td.expected {
					t.Errorf("[OutputLevel] result was incorrect | Want: %v, Got: %v", td.expected, got)
				}
			},
		)
	}
}

func TestPrint(t *testing.T) {
	tests := []struct {
		testName  string
		data      string
		attribute color.Attribute
	}{
		{
			testName:  "should show note properly",
			data:      "test data",
			attribute: color.FgBlack,
		},
	}

	for _, td := range tests {
		t.Run(td.testName, func(t *testing.T) {
			pkg.Print(td.data, td.attribute)
		})
	}
}

func TestPrintNote(t *testing.T) {
	tests := []struct {
		testName string
		note     models.Note
	}{
		{
			testName: "should show note properly",
			note:     models.Note{}, // Empty note
		},
		{
			testName: "should show note properly",
			note:     models.Note{Body: "Non empty note", Path: map[string]string{services.LOCAL.ToStr(): "non-empty-path"}},
		},
	}

	for _, td := range tests {
		t.Run(td.testName, func(t *testing.T) {
			pkg.PrintNote(td.note)
		})
	}
}

func TestPrintNodes(t *testing.T) {
	tests := []struct {
		testName string
		list     []models.Node
	}{
		{
			testName: "should break function",
			list:     []models.Node{},
		},
		{
			testName: "should show note properly",
			list: []models.Node{
				{Pretty: []string{"icon", "Test TITLE"}},
			},
		},
	}

	for _, td := range tests {
		t.Run(td.testName, func(t *testing.T) {
			pkg.PrintNodes(td.list)
		})
	}
}

func TestPrintSettings(t *testing.T) {
	tests := []struct {
		testName string
		settings models.Settings
	}{
		{
			testName: "should show note properly",
			settings: models.Settings{},
		},
	}

	for _, td := range tests {
		t.Run(td.testName, func(t *testing.T) {
			pkg.PrintSettings(td.settings)
		})
	}
}

func TestPrintErrors(t *testing.T) {
	tests := []struct {
		act  string
		errs []error
	}{
		{
			act: "fetch",
			errs: []error{
				errors.New("new mockerror"),
			},
		},
	}

	for _, td := range tests {
		pkg.PrintErrors(td.act, td.errs)
	}
}

func TestSpinner(t *testing.T) {
	got := pkg.Spinner()

	tests := []struct {
		expected *spinner.Spinner
	}{
		{expected: got},
	}

	for _, td := range tests {
		if td.expected != got {
			t.Errorf("Sum of Spinner was different, Want: %v, Got: %v", td.expected, got)
		}
	}
}
