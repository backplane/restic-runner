package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

const defaultUserConfigFile = "~/.restic-runner.yml"
const defaultUserPidFile = "~/.restic-runner.pid"
const defaultSystemConfigFile = "/etc/restic-runner.yml"
const defaultSystemPidFile = "/var/run/restic-runner.pid"

const cmdNameCycle = "cycle"
const cmdNameCommand = "command"

var (
	// version, commit, date, builtBy are provided by goreleaser during build
	progname string = "restic-runner"
	version  string = "dev"
	commit   string = "dev"
	date     string = "unknown"
	builtBy  string = "unknown"

	defaultConfigFile string = defaultSystemConfigFile
	defaultPidFile    string = defaultSystemPidFile

	logger *slog.Logger
)

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version %s; commit %s; built on %s; by %s\n", progname, version, commit, date, builtBy)
	}
	if os.Geteuid() != 0 {
		defaultConfigFile = defaultUserConfigFile
		defaultPidFile = defaultUserPidFile
	}

	logger = slog.Default()
}

func main() {
	if expanded, err := ExpandTilde(defaultConfigFile); err == nil {
		defaultConfigFile = expanded
	}
	if expanded, err := ExpandTilde(defaultPidFile); err == nil {
		defaultPidFile = expanded
	}

	app := &cli.App{
		Name:    progname,
		Version: version,
		Usage:   "operate restic with a config file",
		Before: func(ctx *cli.Context) error {
			// set the log level
			logLevelStr := ctx.String("loglevel")
			switch strings.ToUpper(logLevelStr) {
			case "INFO":
				slog.SetLogLoggerLevel(slog.LevelInfo)
			case "WARN":
				slog.SetLogLoggerLevel(slog.LevelWarn)
			case "ERROR":
				slog.SetLogLoggerLevel(slog.LevelError)
			case "DEBUG":
				slog.SetLogLoggerLevel(slog.LevelDebug)
			default:
				return fmt.Errorf("FATAL: unable to parse loglevel value: '%s'", logLevelStr)
			}
			logger.Debug("starting up",
				"version", version,
				"commit", commit,
				"date", date,
				"builder", builtBy,
			)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:   cmdNameCycle,
				Usage:  "run a restic backup / prune / check cycle using values in a config file",
				Action: MegaHandler,
			},
			{
				Name:   cmdNameCommand,
				Usage:  "run arbitrary restic commands using values in a config file",
				Action: MegaHandler,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: defaultConfigFile,
				Usage: "path to config file",
			},
			&cli.StringFlag{
				Name:  "loglevel",
				Value: "INFO",
				Usage: "how verbosely to log, one of: DEBUG, INFO, WARN, ERROR",
			},
			&cli.StringFlag{
				Name:  "pidfile",
				Value: defaultPidFile,
				Usage: "path to pid lock file; this file prevents issues concurrent jobs",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Error("execution failed", "error", err)
		os.Exit(1)
	}
}
