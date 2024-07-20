package main

import (
	"fmt"

	"github.com/jinzhu/configor"
	"github.com/urfave/cli/v2"
)

func MegaHandler(ctx *cli.Context) error {
	pidfile, err := MakePIDFile(ctx.String("pidfile"))
	if err != nil {
		return fmt.Errorf("FATAL: failed to write pidfile; error:%s", err)
	}
	defer func() {
		if err := pidfile.Close(); err != nil {
			logger.Error("FATAL: failed to remove pidfile", "error", err)
		}
	}()

	r := ResticConfig{}
	if err := configor.Load(&r, ctx.String("config")); err != nil {
		return fmt.Errorf("FATAL: failed to load config; error:%s", err)
	}
	logger.Debug("starting with config", "config", r)

	var result error
	switch ctx.Command.Name {
	case cmdNameCycle:
		result = r.Cycle()
	case cmdNameCommand:
		result = r.Command(ctx.Args().Slice())
	default:
		return fmt.Errorf("unknown command %s", ctx.Command.Name)
	}

	return result
}

// Cycle is a cli handler for running backups
func (r *ResticConfig) Cycle() error {
	if err := r.config_check(); err != nil {
		logger.Warn("config check failed, possibly repo init is needed, trying that...")
		if err := r.init(); err != nil {
			return fmt.Errorf("FATAL: repo init failed; error:%s", err)
		}
		logger.Info("repo init complete")
	}

	logger.Info("starting backup")
	if err := r.backup(); err != nil {
		return fmt.Errorf("FATAL: failed to backup; error:%s", err)
	}
	logger.Info("backup complete")

	logger.Info("checking backups")
	if err := r.backup_check(); err != nil {
		return fmt.Errorf("FATAL: failed to check backups; error:%s", err)
	}
	logger.Info("check complete")

	logger.Info("cleaning up old backups")
	if err := r.forget(); err != nil {
		return fmt.Errorf("FATAL: failed to cleanup old backups; error:%s", err)
	}
	logger.Info("clean-up complete")

	return nil
}

// Command is a cli handler for running arbitrary restic commands
func (r *ResticConfig) Command(args []string) error {
	logger.Debug("running restic command", "command", args)

	if err := r.command(args); err != nil {
		return fmt.Errorf("FATAL: failed to run command %s; error:%s", args, err)
	}
	logger.Debug("command complete")

	return nil
}
