package hareply

import (
	"io/fs"
	"log/slog"
	"net"
	"sync"
)

// App is a hareply application.
// An application reads responses from a file and serves them over TCP.
// It can only be started once.
type App struct {
	sync.Mutex
	host string
	port int

	// fs is the file system to use for reading the response file. If `nil`, the host file system is used.
	fs                fs.FS
	filepath          string
	lastValidResponse []byte

	logger   *slog.Logger
	listener net.Listener
}

// Option is a configuration option for the hareply application.
type Option func(*App)

// New creates a new hareply application that reads responses from the given file.
func New(filepath string, opts ...Option) (*App, error) {
	app := &App{
		filepath: filepath,
		logger:   slog.Default(),
	}

	for _, opt := range opts {
		opt(app)
	}
	return app, nil
}

// WithPort sets the port to listen on for the hareply application.
// If not set, it will pick an available port.
func WithPort(port int) Option {
	return func(a *App) {
		a.port = port
	}
}

// WithFS sets the file system to use for the hareply application.
// It defaults to the host file system.
func WithFS(fs fs.FS) Option {
	return func(a *App) {
		a.fs = fs
	}
}

// WithHost sets the host to listen on for the hareply application.
// It defaults to an empty string `""` which means "all interfaces".
func WithHost(host string) Option {
	return func(a *App) {
		a.host = host
	}
}

// WithLogger sets the logger to use for the hareply application.
// It defaults to `slog.Default()“.
func WithLogger(logger *slog.Logger) Option {
	return func(a *App) {
		a.logger = logger
	}
}

// ListenAddr returns the address the hareply application is listening on.
// This will return `nil` if the application is not currently serving.
func (a *App) ListenAddr() net.Addr {
	a.Lock()
	defer a.Unlock()

	if a.listener == nil {
		return nil
	}
	return a.listener.Addr()
}
