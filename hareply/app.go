package hareply

import (
	"io/fs"
	"net"
	"sync"
)

// App is a hareply application.
type App struct {
	sync.Mutex
	host string
	port int

	// fs is the file system to use for reading the response file. If `nil`, the host file system is used.
	fs       fs.FS
	filepath string
	response []byte

	listener net.Listener
}

// Option is a configuration option for the hareply application.
type Option func(*App)

// New creates a new hareply application that reads responses from the given file.
func New(filepath string, opts ...Option) (*App, error) {
	app := &App{
		filepath: filepath,
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
// It defaults to """
func WithHost(host string) Option {
	return func(a *App) {
		a.host = host
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
