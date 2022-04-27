package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/leoluz/gitops-actions/internal/git"
)

type Config struct {
	BaseSHA    string        `required:"true" split_words:"true"`
	EventSHA   string        `required:"true" split_words:"true"`
	RepoURL    string        `required:"true" split_words:"true"`
	CloneDir   string        `default:"/app/repo" split_words:"true"`
	CmdTimeout time.Duration `default:"40s" split_words:"true"`
}

const (
	EnvPrefix = "GOA"
)

func main() {
	var c Config
	err := envconfig.Process(EnvPrefix, &c)
	if err != nil {
		log.Fatalf("error loading env vars: %s", err)
	}

	gitVersion, err := git.Version(c.CmdTimeout)
	if err != nil {
		log.Fatalf("error checking git version: %s", err)
	}
	log.Print(gitVersion)
	err = git.Clone(c.RepoURL, c.CloneDir, c.CmdTimeout)
	if err != nil {
		log.Fatalf("error cloning repo %s: %s", c.RepoURL, err)
	}

	files, err := git.NewFiles(c.CloneDir, c.BaseSHA, c.EventSHA, c.CmdTimeout)
	if err != nil {
		log.Fatalf("error checking for new files: %s", err)
	}
	for _, file := range files {
		log.Print(file)
	}
}
