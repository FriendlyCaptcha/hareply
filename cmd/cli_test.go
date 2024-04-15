package cmd

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/require"

	"github.com/friendlycaptcha/hareply/testhelper"
)

func TestCLI(t *testing.T) {
	t.Parallel()

	t.Run("serve happy", func(t *testing.T) {
		t.Parallel()
		f := testhelper.SetupAgentStateFile(t, "up 50%\n")
		port := "8302" // Arbitrary port, we could also use 0 and read it from the stdout of the CLI command.

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		go func() {
			out := &bytes.Buffer{}
			err := CLI(ctx, out, []string{"serve", "-f", f, "-p", port}, kong.Exit(func(c int) {
				if c != 0 {
					require.Fail(t, "unexpected exit code", "exit code: %d", c)
				}
			}))
			require.NoError(t, err)
		}()

		time.Sleep(50 * time.Millisecond)
		resp := testhelper.DialAndGetResponse(t, "localhost:"+port)
		require.Equal(t, "up 50%\n", resp)
	})

	t.Run("unknown command", func(t *testing.T) {
		t.Parallel()

		err := CLI(context.Background(), nil, []string{"unknown-command"})
		require.Error(t, err)
	})
}
