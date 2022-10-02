//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package assets

import "fmt"

// ShortSlog is a main slogan of notya.
var ShortSlog = "📝 Take notes quickly and expeditiously"

// MinimalisticBanner is a first banner of notya.
var MinimalisticBanner = `
  _   _    ___    _____  __   __    _    
 | \ | |  / _ \  |_   _| \ \ / /   / \   
 |  \| | | | | |   | |    \   /   / _ \  
 | |\  | | |_| |   | |     | |   / ___ \ 
 |_| \_|  \___/    |_|     |_|  /_/   \_\
`

// GenerateBanner merges slog and banner together, to get final result of application banner.ƒ
func GenerateBanner(banner string, slog string) string {
	template := `
  %v
 %v
   `

	return fmt.Sprintf(template, banner, slog)
}
