package main

import (
	"log"
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

}

func main() {
	app := &cli.App{
		Name:  "restic-runner",
		Usage: "run restic backups from a config file",
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
		},
		Action: func(ctx *cli.Context) error {
			setLogLevel(ctx.String("loglevel"))
			logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel}))

			logger.Debug("starting up",
				"version", version,
				"commit", commit,
				"date", date,
				"builder", builtBy,
			)

			conf := &ResticConfig{}
			if err := configor.Load(conf, ctx.String("config")); err != nil {
				log.Fatalf("failed to load config; error:%+v", err)
			}
			logger.Debug("starting with config", "config", conf)

			if err := conf.check(); err != nil {
				logger.Info("config check failed, possibly repo init is needed, trying that...")
				if err := conf.init(); err != nil {
					log.Fatalf("repo init failed; error:%s", err)
				}
				logger.Info("repo init complete")
			}

			logger.Info("starting backup")
			if err := conf.backup(); err != nil {
				log.Fatalf("failed to backup; error:%+v", err)
			}

			logger.Info("cleaning up old backups")
			if err := conf.forget(); err != nil {
				log.Fatalf("failed to cleanup old backups; error:%+v", err)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
