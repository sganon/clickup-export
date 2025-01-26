package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/urfave/cli/v2"
)

const exitFail = 1

func run(ctx context.Context) error {
	app := &cli.App{
		Name:  "clickup-export",
		Usage: "A CLI tool to export data from ClickUp",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "token",
				Usage:   "Your personal API token",
				EnvVars: []string{"CLICKUP_API_TOKEN"},
			},
			&cli.StringFlag{
				Name:  "workspace-id",
				Usage: "Your workspace ID",
			},
		},
		Action: func(c *cli.Context) error {
			cup := NewClickupClient(c.String("token"))
			wks := WorkspaceID(c.String("workspace-id"))

			docs, err := cup.SearchDocs(ctx, wks)
			if err != nil {
				return fmt.Errorf("error searching docs: %w", err)
			}

			fmt.Println(docs)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		return fmt.Errorf("error running app: %w", err)
	}

	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}
