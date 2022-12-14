//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package models

import (
	"encoding/json"
)

// Note is the main note model of application.
//
//	Example:
//
// ╭─────────────────────────────────────────────╮
// │ Title: new_note.txt                         │
// │ Path: /User/random-user/notya/new_note.txt  │
// │ Body: ... Note content here ...             │
// ╰─────────────────────────────────────────────╯
type Note struct {
	Title string            `json:"title"`
	Path  map[string]string `json:"path"`
	Body  string            `json:"body"`
}

// GetPath returns exact path of provided service.
// If path for provided service doesn't exists result will be empty string.
func (n *Note) GetPath(service string) string {
	return n.Path[service]
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
	return Node{Type: FILE, Title: n.Title, Path: n.Path, Body: n.Body}
}
