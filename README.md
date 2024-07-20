# restic-runner

A simple CLI tool for calling restic with configurations described in a yaml config file.

## Usage

The program emits the following help text when invoked with the '-h' or '--help' flags:

```
NAME:
   restic-runner - run restic from a config file

USAGE:
   restic-runner [global options] command [command options]

VERSION:
   v0.5.0

COMMANDS:
   cycle    run a restic backup / prune / check cycle using values in a config file
   command  run arbitrary restic commands using values in a config file
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value    path to config file (default: "/Users/user/.restic-runner.yml")
   --loglevel value  how verbosely to log, one of: DEBUG, INFO, WARN, ERROR (default: "INFO")
   --pidfile value   path to pid lock file; this file prevents issues concurrent jobs (default: "/Users/user/.restic-runner.pid")
   --help, -h        show help
   --version, -v     print the version
```

## Config File format

Here's an example config file

```yaml
repo: "s3:s3.aws-blerg.example.com:/my-fancy-bucket/restic"
password: "your-restic-password-goes-here"
env:
  AWS_ACCESS_KEY_ID: 0000000000000000000000000
  AWS_SECRET_ACCESS_KEY: 0000000000000000000000000000000
backup_args:
  - "/etc"
  - "/home"
  - "/root"
```
