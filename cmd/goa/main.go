package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/leoluz/gitops-actions/internal/action"
	"github.com/leoluz/gitops-actions/internal/git"
)

type Config struct {
	BaseSHA                  string        `required:"true" split_words:"true"`
	EventSHA                 string        `required:"true" split_words:"true"`
	EventRefName             string        `required:"true" split_words:"true"`
	RepoURL                  string        `required:"true" split_words:"true"`
	ActionDir                string        `default:"go-actions" split_words:"true"`
	CloneDir                 string        `default:"/app/repo" split_words:"true"`
	CmdTimeout               time.Duration `default:"40s" split_words:"true"`
	TwitterConsumerKey       string        `split_words:"true"`
	TwitterConsumerSecret    string        `split_words:"true"`
	TwitterAccessToken       string        `split_words:"true"`
	TwitterAccessTokenSecret string        `split_words:"true"`
}

const (
	EnvPrefix   = "GOA"
	TweetAction = "tweet"
)

func main() {
	cfg := readConfig()
	reg := initRegistry(cfg)
	log.Printf("Starting gitops-actions in %s", cfg.RepoURL)
	log.Printf("base commit SHA: %s", cfg.BaseSHA)
	log.Printf("event commit SHA: %s", cfg.EventSHA)

	actions, err := buildActions(reg, cfg)
	if err != nil {
		log.Fatalf("error building actions: %s", err)
	}
	log.Println("starting executing actions")
	err = run(actions)
	if err != nil {
		log.Fatalf("error running actions: %s", err)
	}
	log.Printf("total actions executed: %d", len(actions))
}

func readConfig() Config {
	var c Config
	err := envconfig.Process(EnvPrefix, &c)
	if err != nil {
		log.Fatalf("error loading env vars: %s", err)
	}
	return c
}

func buildActions(reg *action.Registry, cfg Config) ([]action.Action, error) {
	files, err := getGitFiles(cfg)
	if err != nil {
		return nil, err
	}
	return action.BuildActions(reg, cfg.ActionDir, files)
}

func initRegistry(c Config) *action.Registry {
	tc := action.NewTwitterConfig(c.TwitterConsumerKey, c.TwitterConsumerSecret, c.TwitterAccessToken, c.TwitterAccessTokenSecret)
	registry := action.NewRegistry()
	registry.Add(TweetAction, tc)
	return registry
}

func getGitFiles(c Config) ([]*git.File, error) {
	gitVersion, err := git.Version(c.CmdTimeout)
	if err != nil {
		return nil, fmt.Errorf("error checking git version: %s", err)
	}
	log.Print(gitVersion)
	err = git.Clone(c.RepoURL, c.CloneDir, c.CmdTimeout)
	if err != nil {
		return nil, fmt.Errorf("error cloning repo %s: %s", c.RepoURL, err)
	}

	err = git.Checkout(c.CloneDir, c.EventRefName, c.CmdTimeout)
	if err != nil {
		return nil, fmt.Errorf("error checking out refName %s: %s", c.EventRefName, err)
	}

	files, err := git.GetFiles(c.CloneDir, c.BaseSHA, c.EventSHA, c.CmdTimeout)
	if err != nil {
		return nil, fmt.Errorf("error getting git files: %s", err)
	}
	return files, nil
}

func run(actions []action.Action) error {
	for _, action := range actions {
		err := action.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}
