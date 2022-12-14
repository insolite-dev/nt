//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package models

// Folder is a w-directory representation data structure.
//
//	EXAMPLE:
//
// ╭───────────────╮
// │ ~/notya-path/ │
// │───────────────╯
// │─ new_note.txt
// │─ todo/   ◀──── Sub directory of main notes folder.
// │  │── today.md
// │  │── tomorrow.md
// │  ╰── insolite-notya/  ◀── Sub directory of "todo" folder.
// │      │── issues.txt
// │      ╰── features.txt
// │─ ted-talks.tx
// ╰─ psyco.txt
type Folder struct {
	// Title is the name(not path) of "current" folder.
	Title string `json:"title"`

	// Path is the full-path string name of "current" folder.
	Path map[string]string `json:"path"`
}

// GetPath returns exact path of provided service.
// If path for provided service doesn't exists result will be empty string.
func (f *Folder) GetPath(service string) string {
	return f.Path[service]
}

// ToNode converts [Folder] model to [Node] model.
func (n *Folder) ToNode() Node {
	return Node{Type: FOLDER, Title: n.Title, Path: n.Path}
}
