package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type ForgetConfig struct {
	Days   int `default:"5"`
	Weeks  int `default:"4"`
	Months int `default:"3"`
	Years  int `default:"2"`
}

type ResticConfig struct {
	Repo       string
	Password   string
	Env        map[string]string
	BackupArgs []string `yaml:"backup_args"`
	Forget     ForgetConfig
}

func (r *ResticConfig) cmd(args ...string) *exec.Cmd {
	logger.Debug("restic call", "args", args)
	cmd := exec.Command("restic", args...)
	cmd.Env = cmd.Environ()
	cmd.Env = append(
		cmd.Env,
		fmt.Sprintf("%s=%s", "RESTIC_REPOSITORY", r.Repo),
		fmt.Sprintf("%s=%s", "RESTIC_PASSWORD", r.Password),
	)
	for k, v := range r.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// let the child processes have our STDIO
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func (r *ResticConfig) init() error {
	cmd := r.cmd("init")
	return cmd.Run()
}

func (r *ResticConfig) backup_check() error {
	cmd := r.cmd("check")
	return cmd.Run()
}

func (r *ResticConfig) config_check() error {
	cmd := r.cmd("cat", "config")
	return cmd.Run()
}

func (r *ResticConfig) backup() error {
	backupArgs := []string{"backup", "--verbose"}
	backupArgs = append(backupArgs, r.BackupArgs...)
	cmd := r.cmd(backupArgs...)
	return cmd.Run()
}

func (r *ResticConfig) command(args []string) error {
	cmd := r.cmd(args...)
	return cmd.Run()
}

func (r *ResticConfig) forget() error {
	cmd := r.cmd(
		"forget",
		"-d", strconv.Itoa(r.Forget.Days),
		"-w", strconv.Itoa(r.Forget.Weeks),
		"-m", strconv.Itoa(r.Forget.Months),
		"-y", strconv.Itoa(r.Forget.Years),
		"--prune",
		"--compact",
	)
	return cmd.Run()
}
