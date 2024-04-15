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

	slog.DebugContext(ctx, "initial response loaded", slog.String("response", string(response)))

	listener, err := net.Listen("tcp", net.JoinHostPort(a.host, fmt.Sprint(a.port)))
	if err != nil {

		return fmt.Errorf("failed to listen: %w", err)
	}
	defer listener.Close()

	a.Lock()
	a.listener = listener
	a.Unlock()

	slog.InfoContext(ctx, "listening...", slog.String("address", listener.Addr().String()))

	go func() {
		for {
			select {
			case <-ctx.Done():
				listener.Close()
				slog.InfoContext(ctx, "listener closed")
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					if errors.Is(ctx.Err(), context.Canceled) {
						continue
					}
					slog.ErrorContext(ctx, "failed to accept tcp connection", slog.String("error", err.Error()))
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

	slog.DebugContext(ctx, "updating response")

	response, err := a.updateResponse()
	if err != nil {
		slog.ErrorContext(ctx, "failed to update response, will use cached value", slog.String("error", err.Error()))
	}

	slog.DebugContext(ctx, "writing response", slog.String("response", string(response)))

	_, err = conn.Write(response)
	if err != nil {
		slog.ErrorContext(ctx, "failed to write response", slog.String("error", err.Error()))
		return
	}

	slog.DebugContext(ctx, "response written, closing connection", slog.String("response", string(response)))
}
