package command

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"
)

type Results struct {
	StdOut string
	StdErr string
}

func newResults(stdout, stderr string) Results {
	return Results{
		StdOut: stdout,
		StdErr: stderr,
	}
}

func Run(cmd *exec.Cmd, timeout time.Duration) (Results, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Printf("running command: %s", cmd.String())
	err := cmd.Start()
	if err != nil {
		return newResults(stdout.String(), stderr.String()), err
	}

	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	var timoutCh <-chan time.Time
	if timeout != 0 {
		timoutCh = time.NewTimer(timeout).C
	}

	select {
	case <-timoutCh:
		_ = cmd.Process.Kill()
		err = fmt.Errorf("%s timeout elapsed", timeout.String())
	case e := <-done:
		err = e
	}

	return newResults(stdout.String(), stderr.String()), err
}
