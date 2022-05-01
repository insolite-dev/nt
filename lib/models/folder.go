// Copyright 2022-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package models

// Folder is a w-directory representation data structure.
//
//  EXAMPLE:
// ╭───────────────╮
// │ ~/notya-path/ │
// │───────────────╯
// │─ new_note.txt
// │─ todo/   ◀──── Sub directory of main notes folder.
// │  │── today.md
// │  │── tomorrow.md
// │  ╰── anon-notya/  ◀── Sub directory of "todo" folder.
// │      │── issues.txt
// │      ╰── features.txt
// │─ ted-talks.tx
// ╰─ psycology_resources.txt
//
type Folder struct {
	// Title is the name(not path) of "current" folder.
	Title string `json:"title"`

	// Path is the full-path string name of "current" folder.
	Path string `json:"path"`

	// Files is the slice of "current" folder's sub-nodes.
	// Includes full-path strings of the files/folders names.
	Files []string `json:"files"`
}

// ToNode converts [Folder] model to [Node] model.
func (n *Folder) ToNode() Node {
	return Node{
		Title: n.Title, Path: n.Path,
		Pretty: []string{"", n.Title},
	}
}
