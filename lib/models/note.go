// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package models

// Note is main note model of application.
type Note struct {
	Title string `json:"title"`
	Path  string `json:"path"`
	Body  []byte `json:"body"`
}
