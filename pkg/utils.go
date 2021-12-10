// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package pkg

import (
	"io/ioutil"
	"os"
)

// NotyaPWD, generates path of notya's notes directory.
func NotyaPWD() (*string, error) {
	// Take current user's home directory.
	uhd, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Add notes path
	path := uhd + "/" + "notya/"

	return &path, nil
}

// FileExists, checks if any type of file exists at given path.
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// NewFile, creates new file and writes to its data.
func NewFile(path string, body []byte) error {
	err := os.WriteFile(path, body, 0o600)
	if err != nil {
		return err
	}

	return nil
}

// NewFolder, creates new empty working directory.
func NewFolder(name string) error {
	if err := os.Mkdir(name, 0o750); err != nil {
		return err
	}

	return nil
}

// Delete, removes file or folder, from given path.
func Delete(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

// ReadBody, opens file from given path, and takes its body to return.
func ReadBody(path string) (*string, error) {
	resbyte, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	res := string(resbyte)
	return &res, nil
}

// ListDir, reads all files from given-path directory.
func ListDir(path string) ([]string, error) {
	// Read directory's files.
	list, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Convert list to string list.
	res := []string{}
	for _, d := range list {
		res = append(res, d.Name())
	}

	return res, nil
}
