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

// updateResponse reads the response from the file and updates the `lastValidResponse` field.
// If the response is invalid or could not be read, nil is returned
func (a *App) updateResponse() ([]byte, error) {
	a.Lock()
	defer a.Unlock()

	b, err := a.readResponseFromFile()
	if err != nil {
		return nil, err
	}

	err = ValidateAgentCheckResponse(string(b))
	if err != nil {
		return nil, fmt.Errorf("invalid response in file \"%s\", %w", b, err)
	}

	a.lastValidResponse = b
	return b, nil
}
