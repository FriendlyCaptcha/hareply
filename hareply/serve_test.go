package hareply

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/friendlycaptcha/hareply/testhelper"
)

// startApp starts the application and returns the address it is listening on.
func startApp(ctx context.Context, t *testing.T, app *App) string {
	t.Helper()

	go func() {
		err := app.Serve(ctx)
		require.NoError(t, err)
	}()

	// TODO: we should wait for the server to be ready
	// but we don't have a way to know when it is ready without more code.
	time.Sleep(100 * time.Millisecond)

	return app.ListenAddr().String()
}

func TestServe(t *testing.T) {
	t.Parallel()
	t.Run("happy", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		file := testhelper.SetupAgentStateFile(t, "test")
		app, err := New(file, WithPort(0))
		require.NoError(t, err)

		addr := startApp(ctx, t, app)
		assert.Equal(t, "test", testhelper.DialAndGetResponse(t, addr))

		// Change contents of the file
		require.NoError(t, os.WriteFile(file, []byte("new"), 0644))

		assert.Equal(t, "new", testhelper.DialAndGetResponse(t, addr))

		// Delete the file
		require.NoError(t, os.Remove(file))

		// The response should still be the last one
		assert.Equal(t, "new", testhelper.DialAndGetResponse(t, addr))

	})

	t.Run("already serving", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		file := testhelper.SetupAgentStateFile(t, "test")
		app, err := New(file, WithPort(0))
		require.NoError(t, err)

		startApp(ctx, t, app)
		require.ErrorContains(t, app.Serve(ctx), "already serving")
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		file := "nonexistent"
		app, err := New(file, WithPort(0))
		require.NoError(t, err)

		require.ErrorIs(t, app.Serve(ctx), os.ErrNotExist)
	})

	t.Run("port in use", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		file := testhelper.SetupAgentStateFile(t, "test")
		app, err := New(file, WithPort(0))
		require.NoError(t, err)

		addr := startApp(ctx, t, app)

		_, port, err := net.SplitHostPort(addr)
		require.NoError(t, err)

		portUint, err := net.DefaultResolver.LookupPort(ctx, "tcp", port)
		require.NoError(t, err)

		// Start another app on the same port
		app2, err := New(file, WithPort(portUint))
		require.NoError(t, err)

		require.Error(t, app2.Serve(ctx))
	})

}
