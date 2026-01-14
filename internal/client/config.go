package client

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/mizuchilabs/mantraed/internal/util"
	"github.com/urfave/cli/v3"
)

type Config struct {
	Token               string
	ServerURL           string
	ProfileID           int64
	AgentID             string
	ActiveIP            string
	HealthCheckInterval time.Duration
	UpdateInterval      time.Duration
	ConnectionTimeout   time.Duration
	HealthTimeout       time.Duration
}

var once sync.Once

func Load(cmd *cli.Command) (*Config, error) {
	config := &Config{}
	var err error

	once.Do(func() {
		if cmd == nil {
			return
		}

		initLogger(cmd)

		token := cmd.String("token")
		host := cmd.String("host")
		if token == "" || host == "" {
			err = errors.New("token and host are required")
			return
		}

		profileID, agentID, errTok := parseToken(token)
		if errTok != nil {
			err = errTok
			return
		}

		config = &Config{
			Token:               token,
			ServerURL:           util.CleanURL(host),
			ProfileID:           profileID,
			AgentID:             agentID,
			ActiveIP:            "",
			HealthCheckInterval: 15 * time.Second,
			UpdateInterval:      10 * time.Second,
			ConnectionTimeout:   10 * time.Second,
			HealthTimeout:       5 * time.Second,
		}
	})
	return config, err
}

func initLogger(cmd *cli.Command) {
	level := slog.LevelInfo
	if cmd.Bool("debug") {
		level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(
		tint.NewHandler(colorable.NewColorable(os.Stderr), &tint.Options{
			Level:      level,
			TimeFormat: time.RFC3339,
			NoColor:    !isatty.IsTerminal(os.Stderr.Fd()),
		}),
	))
}

func parseToken(token string) (int64, string, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return 0, "", errors.New("invalid token format")
	}

	profileID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", errors.New("invalid profile ID in token")
	}

	return profileID, parts[1], nil
}
