//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package models

var (
	NotePretty   = ""
	FolderPretty = ""
)

// Node is general purpose data object that used as abstract of
// [Folder] and [Note] structure models.
type Node struct {
	// Title is the name(not path) of "current" node.
	Title string `json:"title"`

	// Path is the full-path string name of "current" node.
	Path string `json:"path"`

	// Pretty is Title but powered with ascii emojis.
	Pretty []string `json:"pretty"`

	// A field representation of [Note]'s [Body].
	Body string `json:"body"`
}

// EditNote is wrapper structure used to
// store two [new/current] nodes.
type EditNode struct {
	Current Node `json:"current"`
	New     Node `json:"new"`
}

// ToNote converts [Node] object to [Note].
func (n *Node) ToNote() Note {
	return Note{Title: n.Title, Path: n.Path, Body: n.Body}
}

// ToFile converts [Node] object to [Folder].
func (n *Node) ToFolder() Folder {
	return Folder{Title: n.Title, Path: n.Path}
}

// StructAsFolder formats the [Node] object as a proper [Folder].
func (n *Node) StructAsFolder() Node {
	var title, path string = n.Title, n.Path

	if len(title) != 0 && string(title[len(title)-1]) != "/" {
		title += "/"
	}

	if len(path) != 0 && string(path[len(path)-1]) != "/" {
		path += "/"
	}

	return Node{Title: title, Path: path, Pretty: n.Pretty}
}

// StructAsNote formats the [Node] object as a proper [Note].
func (n *Node) StructAsNote() Node {
	var title, path string = n.Title, n.Path

	if len(title) != 0 && string(title[len(title)-1]) == "/" {
		title = title[:len(title)-1]
	}

	if len(path) != 0 && string(path[len(path)-1]) == "/" {
		path = path[:len(path)-1]
	}

	return Node{Title: title, Path: path, Pretty: n.Pretty}
}
