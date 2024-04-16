package hareply

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
)

// Serve starts the hareply application with the given arguments.
// It will block until the context is canceled.
func (a *App) Serve(ctx context.Context) error {
	if a.listener != nil {
		return errors.New("already serving")
	}

	response, err := a.updateResponse()
	if err != nil {
		return fmt.Errorf("failed to perform initial update: %w", err)
	}

	a.logger.DebugContext(ctx, "initial response loaded", slog.String("response", string(response)))

	lc := net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp", net.JoinHostPort(a.host, fmt.Sprint(a.port)))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	defer func() {
		closeErr := listener.Close()
		if closeErr != nil {
			a.logger.ErrorContext(ctx, "failed to close listener", slog.String("error", closeErr.Error()))
		}
	}()

	a.Lock()
	a.listener = listener
	a.Unlock()

	a.logger.InfoContext(ctx, "listening...", slog.String("address", listener.Addr().String()))

	go func() {
		for {
			select {
			case <-ctx.Done():
				a.logger.InfoContext(ctx, "listener closing")
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					if errors.Is(ctx.Err(), context.Canceled) {
						continue
					}
					a.logger.ErrorContext(ctx, "failed to accept tcp connection", slog.String("error", err.Error()))
					continue
				}
				go a.handle(ctx, conn)
			}
		}
	}()

	<-ctx.Done()
	return nil
}

func (a *App) handle(ctx context.Context, conn net.Conn) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	a.logger.DebugContext(ctx, "updating response")

	response, err := a.updateResponse()
	if err != nil {
		response = a.lastValidResponse
		a.logger.ErrorContext(ctx, "failed to update response, will use cached value", slog.String("error", err.Error()))
	}

	a.logger.DebugContext(ctx, "writing response", slog.String("response", string(response)))

	_, err = conn.Write(response)
	if err != nil {
		a.logger.ErrorContext(ctx, "failed to write response", slog.String("error", err.Error()))
		return
	}

	a.logger.DebugContext(ctx, "response written, closing connection", slog.String("response", string(response)))
}
