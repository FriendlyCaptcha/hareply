// Package cmd contains the CLI implementation for hareply.
package cmd

import (
	"context"
	"fmt"
	"io"

	"log/slog"

	"github.com/alecthomas/kong"
	"github.com/friendlycaptcha/hareply/buildinfo"
	"github.com/friendlycaptcha/hareply/hareply"
)

// Config represents the configuration for the hareply CLI.
type Config struct {
	Debug bool `help:"Enable debug mode."`

	Serve struct {
		Host string `help:"Host to listen on." default:""`
		Port int    `help:"Port to reply on." short:"p" default:"8442"`
		File string `help:"Path to read response from." short:"f" default:"agentstate"`
	} `cmd:"" help:"Start hareply TCP responder service."`

	Version struct{} `cmd:"" help:"Show version information."`
}

// CLI runs the hareply CLI.
func CLI(ctx context.Context, w io.Writer, args []string, opts ...kong.Option) error {
	cli := Config{}

	opts = append(opts,
		kong.Name("hareply"),
		kong.Description("A simple TCP server that replies with a response read from a file."),
		kong.DefaultEnvars("HAREPLY_"),
	)

	kcli, err := kong.New(&cli, opts...)
	if err != nil {
		return fmt.Errorf("error creating CLI parser: %w", err)
	}

	kctx, err := kcli.Parse(args)
	if err != nil {
		return err
	}

	logger := setupLogger(w, cli.Debug)

	switch kctx.Command() {
	case "serve":
		app, err := hareply.New(cli.Serve.File,
			hareply.WithPort(cli.Serve.Port),
			hareply.WithHost(cli.Serve.Host),
			hareply.WithLogger(logger),
		)
		if err != nil {
			return err
		}
		return app.Serve(ctx)
	case "version":
		fmt.Fprintln(w, buildinfo.FullVersion())
		return nil
	default:
		return fmt.Errorf("unknown command: %s", kctx.Command())
	}
}

func setupLogger(w io.Writer, debug bool) *slog.Logger {
	lvl := slog.LevelInfo
	if debug {
		lvl = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: lvl,
	}))
	return logger
}
