package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/jinzhu/configor"
	"github.com/urfave/cli/v2"
)

type ExpireConfig struct {
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
	Expire     ExpireConfig
}

func (c *ResticConfig) resticCommand(args ...string) *exec.Cmd {
	cmd := exec.Command("restic", args...)
	log.Printf("args: %+v\n", args)
	cmd.Env = cmd.Environ()
	cmd.Env = append(
		cmd.Env,
		fmt.Sprintf("%s=%s", "RESTIC_REPOSITORY", c.Repo),
		fmt.Sprintf("%s=%s", "RESTIC_PASSWORD", c.Password),
	)
	for k, v := range c.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	// Set the standard input, output, and error of the child process
	// to the same as the parent process's standard handles
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func (c *ResticConfig) resticInit() error {
	cmd := c.resticCommand("init")
	return cmd.Run()
}

func (c *ResticConfig) resticCheck() error {
	cmd := c.resticCommand("cat", "config")
	return cmd.Run()
}

func (c *ResticConfig) resticBackup() error {
	backupArgs := []string{"backup", "--verbose"}
	backupArgs = append(backupArgs, c.BackupArgs...)
	cmd := c.resticCommand(backupArgs...)
	return cmd.Run()
}

func (c *ResticConfig) resticForget() error {
	cmd := c.resticCommand(
		"forget",
		"-d", strconv.Itoa(c.Expire.Days),
		"-w", strconv.Itoa(c.Expire.Weeks),
		"-m", strconv.Itoa(c.Expire.Months),
		"-y", strconv.Itoa(c.Expire.Years),
	)
	return cmd.Run()
}

func main() {
	app := &cli.App{
		Name:  "backup",
		Usage: "run restic backups",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "/etc/restic-runner.yml",
				Usage: "path to config file",
			},
		},
		Action: func(c *cli.Context) error {
			conf := &ResticConfig{}
			if err := configor.Load(conf, c.String("config")); err != nil {
				log.Fatalf("failed to load config; error:%+v", err)
			}
			// log.Printf("config: %+v\n", conf)

			if conf.resticCheck() != nil {
				log.Println("cat check failed, attempting repo init")
				if conf.resticInit() != nil {
					log.Fatal("repo init failed")
				}
			}
			if err := conf.resticBackup(); err != nil {
				log.Fatalf("failed to backup; error:%+v", err)
			}
			if err := conf.resticForget(); err != nil {
				log.Fatalf("failed to cleanup old backups; error:%+v", err)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
