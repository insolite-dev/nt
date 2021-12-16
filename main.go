// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/anonistas/notya/cmd"
	"github.com/dimiro1/banner"
)

func init() {
	bannerBytes, _ := ioutil.ReadFile("assets/banner.txt")
	banner.Init(os.Stdout, true, false, strings.NewReader(string(bannerBytes)))
}

func main() {
	cmd.RunApp()
}
