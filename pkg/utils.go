package pkg

import "os"

// NotyaPWD, generates path of notya's notes directory.
func NotyaPWD() (*string, error) {
	// Take current user's home directory.
	uhd, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Add notes path
	path := uhd + "/" + "notya-notes" + "/"

	return &path, nil
}
