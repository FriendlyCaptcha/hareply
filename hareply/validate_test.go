package hareply

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAgentCheckResponse(t *testing.T) {
	testCases := []struct {
		response string
		ok       bool
	}{
		{"", false},
		{"\n", true},
		{"75%\r", true},
		{"%\r", false},
		{"-1%\r", false},
		{"100%\r", true},
		{"1000%\r", false},
		{"maxconn:30\r", true},
		{"ready\r", true},
		{"drain\r", true},
		{"maint\r", true},
		{"down\r", true},
		{"fail\r", true},
		{"stopped\r", true},
		{"up\r", true},
		{"75%#description\r", false},
		{"maxconn:30#description\r", false},
		{"ready#description\r", false},
		{"drain#description\r", false},
		{"maint#description\r", false},
		{"down#description\r", true},
		{"fail#description\r", true},
		{"stopped#description123\r", true},
		{"up#description\r", false},
		{"75% \n", true},
		{"maxconn:30 \n", true},
		{"ready \n", true},
		{"drain \n", true},
		{"maint \n", true},
		{"down \n", true},
		{"fail \n", true},
		{"stopped \n", true},
		{"up \n", true},
		{"up down\n", true},
		{"up 50%\n", true},
		{"20% 30%\n", true},
		{"20% 30% down#description\n", true},
		{"up\t20%\n", true},
		{"up\t20%\t\t\t\t\t\t\tdown#description\t\n", true},
		{"up\t20% down#description,stopped\n", true},
		{"up\ndown\n", false},
		{"up\n\n", true},
		{"up\n\n\r", true},
		{"up\r\n", true},
		{"up", false},
		{"down#a_b\n", true},
		{"down#a-b\n", true},
		{"down#a.b\n", false},
	}

	for _, tc := range testCases {
		t.Run(tc.response, func(t *testing.T) {
			result := ValidateAgentCheckResponse(tc.response)
			if tc.ok {
				assert.NoError(t, result)
			} else {
				assert.Error(t, result)
			}
		})
	}
}
