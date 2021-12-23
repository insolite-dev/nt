// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package pkg_test

import (
	"fmt"
	"testing"

	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
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

func TestShowNote(t *testing.T) {
	tests := []struct {
		testName string
		note     models.Note
	}{
		{
			testName: "should show note properly",
			note:     models.Note{}, // Empty note
		},
	}

	for _, td := range tests {
		t.Run(td.testName, func(t *testing.T) {
			pkg.ShowNote(td.note)
		})
	}
}

func TestShowListOfNotes(t *testing.T) {
	type args struct {
		list  []string
		limit int
	}

	tests := []struct {
		testName string
		arg      args
	}{
		{
			testName: "should show note properly",
			arg: args{
				list:  []string{"1", "2", "3", "4"},
				limit: 2,
			},
		},
	}

	for _, td := range tests {
		t.Run(td.testName, func(t *testing.T) {
			pkg.ShowListOfNotes(td.arg.list, td.arg.limit)
		})
	}
}
