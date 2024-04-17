package testhelper

import (
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	// TestAgentStateFilepath is the path to the test file.
	TestAgentStateFilepath = "agentstate"
)

func SetupAgentStateFile(t *testing.T, content string) string {
	f := filepath.Join(t.TempDir(), TestAgentStateFilepath)
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))
	return f
}

func DialAndGetResponse(t *testing.T, addr string) string {
	t.Helper()

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	defer conn.Close()

	// Read the response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	require.NoError(t, err)
	return string(buf[:n])
}
