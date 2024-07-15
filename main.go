package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jinzhu/configor"
	"github.com/urfave/cli/v2"
)

var (
	// version, commit, date, builtBy are provided by goreleaser during build
	version = "dev"
	commit  = "dev"
	date    = "unknown"
	builtBy = "unknown"

	logLevel *slog.LevelVar
	logger   *slog.Logger
)

func init() {
	logLevel = new(slog.LevelVar)

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("restic-runner version %s; commit %s; built on %s; by %s\n", version, commit, date, builtBy)
	}
}

func main() {
	app := &cli.App{
		Name:    "restic-runner",
		Version: version,
		Usage:   "run restic backups from a config file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "/etc/restic-runner.yml",
				Usage: "path to config file",
			},
			&cli.StringFlag{
				Name:  "loglevel",
				Value: "INFO",
				Usage: "how verbosely to log, one of: DEBUG, INFO, WARN, ERROR",
			},
			&cli.StringFlag{
				Name:  "pidfile",
				Value: "/var/run/restic-runner.pid",
				Usage: "how verbosely to log, one of: DEBUG, INFO, WARN, ERROR",
			},
		},
		Action: func(ctx *cli.Context) error {
			setLogLevel(ctx.String("loglevel"))
			logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel}))

			pidfile, err := MakePIDFile(ctx.String("pidfile"))
			if err != nil {
				return fmt.Errorf("FATAL: failed to write pidfile; error:%s", err)
			}
			defer func() {
				if err := pidfile.Close(); err != nil {
					logger.Error("FATAL: failed to remove pidfile", "error", err)
				}
			}()

			logger.Debug("starting up",
				"version", version,
				"commit", commit,
				"date", date,
				"builder", builtBy,
			)

			conf := &ResticConfig{}
			if err := configor.Load(conf, ctx.String("config")); err != nil {
				return fmt.Errorf("FATAL: failed to load config; error:%s", err)
			}
			logger.Debug("starting with config", "config", conf)

			if err := conf.check(); err != nil {
				logger.Warn("config check failed, possibly repo init is needed, trying that...")
				if err := conf.init(); err != nil {
					return fmt.Errorf("FATAL: repo init failed; error:%s", err)
				}
				logger.Info("repo init complete")
			}

			logger.Info("starting backup")
			if err := conf.backup(); err != nil {
				return fmt.Errorf("FATAL: failed to backup; error:%s", err)
			}

			logger.Info("cleaning up old backups")
			if err := conf.forget(); err != nil {
				return fmt.Errorf("FATAL: failed to cleanup old backups; error:%s", err)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Error("execution failed", "error", err)
		os.Exit(1)
	}
}
