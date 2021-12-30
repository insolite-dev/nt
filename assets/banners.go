package assets

import "fmt"

// ShortSlog is a main slogan of notya.
var ShortSlog = "📝 Take notes quickly and expeditiously"

// MinimalisticBanner is a first banner of notya.
var MinimalisticBanner = `
  _   _   ___  _____ __   __  _    
 | \ | | / _ \|_   _|\ \ / / / \   
 |  \| || | | | | |   \   / / _ \  
 | |\  || |_| | | |    | | / ___ \ 
 |_| \_| \___/  |_|    |_|/_/   \_\
`

func GenerateBanner(banner string, slog string) string {
	template := `
  %v
 %v
   `

	return fmt.Sprintf(template, banner, slog)
}
