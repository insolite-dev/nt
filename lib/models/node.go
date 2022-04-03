// Copyright 2022-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package models

// Node is general purpose data object that used as abstract of
// [Folder] and [Note] structure models.
type Node struct {
	// Title is the name(not path) of "current" node.
	Title string `json:"title"`

	// Path is the full-path string name of "current" node.
	Path string `json:"path"`
}

// ToNote converts [Node] object to [Note].
func (n *Node) ToNote() Note {
	return Note{Title: n.Title, Path: n.Path}
}

// ToFile converts [Node] object to [Folder].
func (n *Node) ToFolder() Folder {
	return Folder{Title: n.Title, Path: n.Path}
}

// EditNote is wrapper structure used to
// store two [new/current] nodes.
type EditNode struct {
	Current Node `json:"current"`
	New     Node `json:"new"`
}
