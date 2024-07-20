package main

import (
	"fmt"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

var TildePathRegex = regexp.MustCompile(`^~(?P<username>[a-z][a-z0-9._-]{0,63})?(?P<remainder>/.*)?$`)

// ExpandTilde expands the tilde in a path to a user's home directory
func ExpandTilde(path string) (string, error) {
	if !strings.HasPrefix(path, `~`) {
		return path, nil
	}

	matches := TildePathRegex.FindStringSubmatch(path)
	if matches == nil {
		return "", fmt.Errorf("path starts with '~' but doesn't conform to known format")
	}

	username := matches[1]
	remainder := matches[2]
	var usr *user.User
	var err error

	if username == "" {
		usr, err = user.Current()
	} else {
		usr, err = user.Lookup(username)
	}
	if err != nil {
		return "", fmt.Errorf("failed user lookup: %v", err)
	}

	homeDir := usr.HomeDir
	return filepath.Join(homeDir, remainder), nil
}
