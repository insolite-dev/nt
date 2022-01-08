package models

// All editors listed.
var (
	// VI related.
	VI     string = "vi"
	Vim    string = "vim"
	NeoVim string = "nvim"
	MacVim string = "mvim"
	GUIVim string = "gvim"

	// VS-Code related.
	VSCode         string = "code"
	VSCodeInsiders string = "code-insiders"
)

// Settings is a main structure model of application settings.
type Settings struct {
	Editor string `json:"editor" default:"vi"`
}

// InitSettings returns default variant of settings structure model.
func InitSettings() Settings {
	return Settings{Editor: VI}
}
