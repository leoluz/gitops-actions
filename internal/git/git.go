package git

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/leoluz/gitops-actions/internal/command"
)

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

func Checkout(cloneDir, refName string, timeout time.Duration) error {
	cmd := exec.Command("git", "checkout", refName)
	cmd.Dir = cloneDir
	res, err := command.Run(cmd, timeout)
	if err != nil {
		return fmt.Errorf("error executing command %q: %s: stderr: %s", cmd.String(), err, res.StdErr)
	}
	return nil
}
