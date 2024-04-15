package hareply

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	// ErrInvalidAgentCheckResponse is returned when the response to an agent check request is invalid.
	ErrInvalidAgentCheckResponse = errors.New("invalid agent check response value")
)
var regexpValidateAgentCheckResponse = regexp.MustCompile(
	`^((\d{1,3}%|maxconn:\d+|ready|drain|maint|up|((fail|stopped|down)(#\w+)?))([ |\t|,]*))*[\n|\r]+$`,
)

// ValidateAgentCheckResponse validates the response, making sure it is a valid response
// to an agent check request. See the docs here
// https://www.haproxy.com/documentation/haproxy-configuration-manual/latest/#5.2-agent-check
func ValidateAgentCheckResponse(v string) error {
	if !regexpValidateAgentCheckResponse.MatchString(v) {
		// A more human-friendly error for empty strings.
		if len(v) == 0 {
			return fmt.Errorf("%w: empty response", ErrInvalidAgentCheckResponse)
		}
		// And for cases in which there is no newline delimiter.
		if v[len(v)-1] != '\n' && v[len(v)-1] != '\r' {
			return fmt.Errorf("%w: missing newline", ErrInvalidAgentCheckResponse)
		}

		return ErrInvalidAgentCheckResponse
	}
	return nil
}
