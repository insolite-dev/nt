package models

// Note is main note model of application.
type Note struct {
	Title string `json:"title"`
	Path  string `json:"path"`
	Body  string `json:"body"`
}
