package hareply

import (
	"fmt"
	"io"
	"io/fs"
	"os"
)

func (a *App) readResponseFromFile() ([]byte, error) {
	var f fs.File
	var err error

	if a.fs == nil {
		f, err = os.Open(a.filepath)
	} else {
		f, err = a.fs.Open(a.filepath)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", a.filepath, err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", a.filepath, err)
	}
	return b, nil
}

// updateResponse reads the response from the file and updates the response field.
// It returns the new response and any error that occurred. If an error occurred,
// the response field is not updated and the last known response is returned.
func (a *App) updateResponse() ([]byte, error) {
	a.Lock()
	defer a.Unlock()

	b, err := a.readResponseFromFile()
	if err != nil {
		return a.response, err
	}
	a.response = b
	return a.response, nil
}
