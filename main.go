package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mizuchilabs/mantraed/internal/client"
	"github.com/mizuchilabs/mantraed/internal/util"
	"github.com/urfave/cli/v3"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	cmd := &cli.Command{
		EnableShellCompletion: true,
		Suggest:               true,
		Name:                  "mantraed",
		Version:               fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, Date),
		Usage:                 "mantrae daemon",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg, err := client.Load(cmd)
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}

			// Start agent
			slog.Info("Agent starting...", "version", Version)
			client.NewAgent(cfg).Run(ctx)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "update",
				Usage: "Check for updates or update mantraed to the latest version",
				Description: `Check if a newer version of mantraed is available.
Use the --install flag to download and install the latest version.

Note: Automatic installation does not work inside Docker containers.`,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					util.Update(Version, cmd.Bool("install"))
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Enable debug logging",
				Sources: cli.EnvVars("DEBUG"),
			},
			&cli.BoolFlag{
				Name:  "install",
				Usage: "Download and install the latest version (used with update command, does not work in Docker)",
				Value: false,
			},
			&cli.StringFlag{
				Name:    "token",
				Usage:   "Mantrae API token",
				Value:   "",
				Sources: cli.EnvVars("TOKEN"),
			},
			&cli.StringFlag{
				Name:    "host",
				Usage:   "Mantrae API host",
				Value:   "",
				Sources: cli.EnvVars("HOST"),
			},
		},
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := cmd.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
