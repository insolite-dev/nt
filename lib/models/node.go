//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package models

import "os"

var (
	// Early defined pretties.
	NotePretty   = ""
	FolderPretty = ""
)

// NodeType is custom string wrapper to represent node's type.
type NodeType string

var (
	FILE   NodeType = "FILE"
	FOLDER NodeType = "FOLDER"
)

// Node is general purpose data object that used as abstract of
// [Folder] and [Note] structure models.
type Node struct {
	// Type represents the exact type of current node.
	// It can be `file` or `folder`.
	Type NodeType `json:"typ"`

	// Title is the name(not path) of "current" node.
	Title string `json:"title"`

	// Path is the full-path string name of "current" node.
	Path map[string]string `json:"path"`

	// A field representation of [Note]'s [Body].
	Body string `json:"body"`

	// Pretty is Title but powered with ascii emojis.
	// Shouldn't used as a production field.
	Pretty []string `json:"pretty,omitempty"`
}

// EditNote is wrapper structure used to
// store two [new/current] nodes.
type EditNode struct {
	Current Node `json:"current"`
	New     Node `json:"new"`
}

// GetPath returns exact path of provided service.
// If path for provided service doesn't exists result will be empty string.
func (n *Node) GetPath(service string) string {
	return n.Path[service]
}

// UpdatePath updates concrete [service]'s path with [path].
func (n *Node) UpdatePath(service, path string) *Node {
	p := map[string]string{service: path}

	for k, v := range n.Path {
		if k != service {
			p[k] = v
		}
	}

	n.Path = p

	return n
}

func (n *Node) IsFolder() bool {
	return n.Type == FOLDER
}

func (n *Node) IsFile() bool {
	return n.Type == FILE
}

// ToNote converts [Node] object to [Note].
func (n *Node) ToNote() Note {
	var title string = n.Title

	if len(title) != 0 && string(title[len(title)-1]) == "/" {
		title = title[:len(title)-1]
	}

	return Note{Title: title, Path: n.Path, Body: n.Body}
}

// ToFile converts [Node] object to [Folder].
func (n *Node) ToFolder() Folder {
	var title string = n.Title

	if len(title) != 0 && string(title[len(title)-1]) != "/" {
		title += "/"
	}

	return Folder{Title: title, Path: n.Path}
}

// PrettyFromEntry generates a pretty icon appropriate to provided entry.
func PrettyFromEntry(e os.DirEntry) string {
	if e.IsDir() {
		return FolderPretty
	}

	return NotePretty
}
