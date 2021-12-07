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

// CreateFile, creates new file and writes to its data.
func CreateFile(path string, body []byte) error {
	err := os.WriteFile(path, body, 0o600)
	if err != nil {
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
