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
			wks := Workspace{ID: WorkspaceID(c.String("workspace-id"))}

			spaces, err := cup.GetWorkspaceSpaces(ctx, wks.ID)
			if err != nil {
				return fmt.Errorf("error getting space from workspace: %w", err)
			}

			for _, s := range spaces {
				fmt.Printf("Space %s:\n", s.Name)
				if err := PopulateSpace(ctx, cup, wks.ID, s); err != nil {
					return fmt.Errorf("error populating space: %w", err)
				}
			}

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
