// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package models

// Note is the main note model of application.
type Note struct {
	Title string `json:"title"`
	Path  string `json:"path"`
	Body  string `json:"body"`
}

// EditNote is a model that has two note fields inside,
// which used to edit note or rename it.
type EditNote struct {
	Current Note `json:"current"`
	New     Note `json:"new"`
}
