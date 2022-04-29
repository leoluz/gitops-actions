package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/leoluz/gitops-actions/internal/command"
)

type GitStatus string

const (
	StatusAdded     GitStatus = "added"
	StatusDeleted   GitStatus = "deleted"
	StatusModified  GitStatus = "modified"
	statusUndefined GitStatus = "undefined"
)

type File struct {
	name     string
	fullPath string
	status   GitStatus
	Content  []byte
}

func (f *File) GetName() string {
	return f.name
}

func (f *File) GetFullpath() string {
	return f.fullPath
}

func (f *File) GetStatus() GitStatus {
	return f.status
}

func (s GitStatus) String() string {
	return string(s)
}

func NewFile(name, cloneDir string, status GitStatus) *File {
	return &File{
		name:     name,
		fullPath: filepath.Join(cloneDir, name),
		status:   status,
	}
}

func Version(timeout time.Duration) (string, error) {
	cmd := exec.Command("git", "version")
	res, err := command.Run(cmd, timeout)
	if err != nil {
		return "", fmt.Errorf("error executing command %q: %s: stderr: %s", cmd.String(), err, res.StdErr)
	}
	return res.StdOut, nil
}

func Clone(repoURL, cloneDir string, timeout time.Duration) error {
	cmd := exec.Command("git", "clone", repoURL, cloneDir)
	res, err := command.Run(cmd, timeout)
	if err != nil {
		return fmt.Errorf("error executing command %q: %s: stderr: %s", cmd.String(), err, res.StdErr)
	}
	return nil
}

func NewFiles(cloneDir, fromSHA, toSHA string, timeout time.Duration) ([]string, error) {
	cmd := exec.Command("git", "diff", "--diff-filter=A", "--name-only", fromSHA+".."+toSHA)
	cmd.Dir = cloneDir
	res, err := command.Run(cmd, timeout)
	if err != nil {
		return nil, fmt.Errorf("error executing command %q: %s: stderr: %s", cmd.String(), err, res.StdErr)
	}
	return strings.Split(res.StdOut, "\n"), nil
}

func GetFiles(cloneDir, fromSHA, toSHA string, timeout time.Duration) ([]*File, error) {
	cmd := exec.Command("git", "diff", "--name-status", fromSHA+".."+toSHA)
	cmd.Dir = cloneDir
	res, err := command.Run(cmd, timeout)
	if err != nil {
		return nil, fmt.Errorf("error executing command %q: %s: stderr: %s", cmd.String(), err, res.StdErr)
	}
	return ToFiles(cloneDir, res.StdOut), nil
}

func ToFiles(cloneDir, diffOutput string) []*File {
	lines := strings.Split(diffOutput, "\n")
	files := []*File{}
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			continue
		}
		status := toStatus(parts[0])
		if status != statusUndefined {
			f := NewFile(parts[1], cloneDir, status)
			files = append(files, f)
		}
	}
	return files
}

func toStatus(code string) GitStatus {
	switch code {
	case "D":
		return StatusDeleted
	case "M":
		return StatusModified
	case "A":
		return StatusAdded
	default:
		return statusUndefined
	}
}

func Checkout(cloneDir, refName string, timeout time.Duration) error {
	cmd := exec.Command("git", "checkout", refName)
	cmd.Dir = cloneDir
	res, err := command.Run(cmd, timeout)
	if err != nil {
		return fmt.Errorf("error executing command %q: %s: stderr: %s", cmd.String(), err, res.StdErr)
	}
	return nil
}
