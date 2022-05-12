// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package models

import "encoding/json"

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

// ToJSON converts string note structure model to map value.
func (s *Note) ToJSON() map[string]interface{} {
	b, _ := json.Marshal(&s)

	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)

	return m
}

// ToNode converts [Note] model to [Node] model.
func (n *Note) ToNode() Node {
	return Node{
		Title: n.Title, Path: n.Path,
		Pretty: []string{NotePretty, n.Title},
	}
}
