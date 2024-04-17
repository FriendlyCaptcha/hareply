// Package main contains the main function for the hareply CLI.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/friendlycaptcha/hareply/cmd"
)

// Run runs the hareply CLI with the given arguments.
func Run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err := cmd.CLI(ctx, w, args)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	if err := Run(ctx, os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
