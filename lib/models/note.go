// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package models

// Note is the main note model of application.
//
//  Example:
// ╭─────────────────────────────────────────────╮
// │ Title: new_note.txt                         │
// │ Path: /User/random-user/notya/new_note.txt  │
// │ Body: ... Note content here ...             │
// ╰─────────────────────────────────────────────╯
type Note struct {
	Title string `json:"title"`
	Path  string `json:"path"`
	Body  string `json:"body"`
}

// EditNote is wrapper structure used to
// store two [new/current] notes.
type EditNote struct {
	Current Note `json:"current"`
	New     Note `json:"new"`
}
