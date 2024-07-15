# restic-runner

A simple CLI tool for calling restic with configurations described in a yaml config file.

## Usage

The program emits the following help text when invoked with the '-h' or '--help' flags:

```
NAME:
   restic-runner - run restic backups from a config file

USAGE:
   restic-runner [global options] command [command options]

VERSION:
   dev

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value    path to config file (default: "/etc/restic-runner.yml")
   --loglevel value  how verbosely to log, one of: DEBUG, INFO, WARN, ERROR (default: "INFO")
   --pidfile value   how verbosely to log, one of: DEBUG, INFO, WARN, ERROR (default: "/var/run/restic-runner.pid")
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
