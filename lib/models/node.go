//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package models

import (
	"encoding/json"
	"os"
	"strings"
)

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
	Body string `json:"body,omitempty"`

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

// RebuildParent updates the parent(s) of node in [Title] and [Path].
func (n *Node) RebuildParent(parentCurrent, parentNew Node, service string, s Settings) *Node {
	if parentCurrent.Title[len(parentCurrent.Title)-1] != '/' {
		parentCurrent.Title += "/"
	}

	if parentNew.Title[len(parentNew.Title)-1] != '/' {
		parentNew.Title += "/"
	}

	split := SplitPath(strings.Replace(n.Title, parentCurrent.Title, parentNew.Title, 1))
	n.Title = CollectPath(split)

	path := "notya"
	if len(s.FirebaseCollection) != 0 {
		path = s.FirebaseCollection
	} else if len(s.Name) != 0 {
		path = s.Name
	}
	if path[len(path)-1] != '/' && n.Title[0] != '/' {
		path += "/"
	}

	return n.UpdatePath(service, path+n.Title)
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

// ToJSON converts node structure model to map value.
func (s *Node) ToJSON() map[string]interface{} {
	b, _ := json.Marshal(&s)

	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)

	return m
}

// FromJson converts provided map data to [Node] structure.
func (s *Node) FromJson(data map[string]interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonBytes, &s)
}

// GenPretty generates default pretty ASCII icon based
// on pointed node's type.
func (s *Node) GenPretty() string {
	if s.IsFolder() {
		return FolderPretty
	}

	return NotePretty
}

// PrettyFromEntry generates a pretty icon appropriate to provided entry.
func PrettyFromEntry(e os.DirEntry) string {
	if e.IsDir() {
		return FolderPretty
	}

	return NotePretty
}

// Split splits the path to fields by char:'/'
func SplitPath(str string) []string {
	split := strings.Split(str, "/")
	for i, s := range split { // remove each empty sub-string.
		if len(s) == 0 {
			split = append(split[:i], split[i+1:]...)
		}
	}

	return split
}

// CollectPath is reversed implementation of SplitPath, which collects
// the fields that is splitted via SplitPath function.
func CollectPath(splitted []string) string {
	res := ""

	for i, s := range splitted {
		res += s
		if res[len(res)-1] != '/' && i != len(splitted)-1 {
			res += "/"
		}
	}

	return res
}
