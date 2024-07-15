package main

import (
	"fmt"
	"os"
)

type PIDFile struct {
	Path       string
	Identifier string
}

// MakePIDFile attempts to create a pid lock file, it will error if the file
// could not be written, which is the case when some other process already has
// a pidfile in place; it is critical to call the Close method when work is
// complete
func MakePIDFile(filePath string) (*PIDFile, error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0444)
	if err != nil {
		return nil, fmt.Errorf("failed to create pidfile; error:%s", err)
	}

	identifier := fmt.Sprintf("%d\n", os.Getpid())
	if _, err := file.WriteString(identifier); err != nil {
		return nil, fmt.Errorf("failed to write to pidfile; error:%s", err)
	}
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("failed to close pidfile after writing; error:%s", err)
	}

	return &PIDFile{Path: filePath, Identifier: identifier}, nil
}

// Close removes the pidfile (as long as it contains our PID)
func (p *PIDFile) Close() error {
	data, err := os.ReadFile(p.Path)
	if err != nil {
		return fmt.Errorf("failed to read pidfile for closing; error:%s", err)
	}

	if fileContent := string(data); fileContent != p.Identifier {
		return fmt.Errorf("unexpected identifier in pid file; have:%s, expect:%s", fileContent, p.Identifier)
	}

	if err := os.Remove(p.Path); err != nil {
		return fmt.Errorf("failed to remove pidfile %s; error:%s", p.Path, err)
	}

	return nil
}
