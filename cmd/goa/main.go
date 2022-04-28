package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/leoluz/gitops-actions/internal/action"
	"github.com/leoluz/gitops-actions/internal/git"
)

type Config struct {
	BaseSHA                  string        `required:"true" split_words:"true"`
	EventSHA                 string        `required:"true" split_words:"true"`
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
	actions, err := buildActions(reg, cfg)
	if err != nil {
		log.Fatalf("error building actions: %s", err)
	}
	err = run(actions)
	if err != nil {
		log.Fatalf("error running actions: %s", err)
	}
}

func readConfig() Config {
	var c Config
	err := envconfig.Process(EnvPrefix, &c)
	if err != nil {
		log.Fatalf("error loading env vars: %s", err)
	}
	return c
}

func initRegistry(c Config) *action.Registry {
	tc := action.NewTwitterConfig(c.TwitterConsumerKey, c.TwitterConsumerSecret, c.TwitterAccessToken, c.TwitterAccessTokenSecret)
	registry := action.NewRegistry()
	registry.Add(TweetAction, tc)
	return registry
}

func buildActions(reg *action.Registry, cfg Config) ([]action.Action, error) {
	newFiles, err := getNewFiles(cfg)
	if err != nil {
		return nil, err
	}

	actions := []action.Action{}
	for _, file := range newFiles {
		if strings.HasPrefix(file, cfg.ActionDir) {
			dirs := strings.Split(file, string(os.PathSeparator))
			if len(dirs) > 2 {
				actionName := dirs[1]
				actionCreator := reg.Get(actionName)
				if actionCreator == nil {
					log.Printf("action %q not found in registry: skipping file %q", actionName, file)
					continue
				}
				content, err := ioutil.ReadFile(file)
				if err != nil {
					return nil, fmt.Errorf("error reading file %q: %s", file, err)
				}
				action := actionCreator.NewAction(string(content))
				actions = append(actions, action)
			}
		}
	}
	return actions, nil
}

func getNewFiles(c Config) ([]string, error) {
	gitVersion, err := git.Version(c.CmdTimeout)
	if err != nil {
		return nil, fmt.Errorf("error checking git version: %s", err)
	}
	log.Print(gitVersion)
	err = git.Clone(c.RepoURL, c.CloneDir, c.CmdTimeout)
	if err != nil {
		return nil, fmt.Errorf("error cloning repo %s: %s", c.RepoURL, err)
	}

	files, err := git.NewFiles(c.CloneDir, c.BaseSHA, c.EventSHA, c.CmdTimeout)
	if err != nil {
		return nil, fmt.Errorf("error checking for new files: %s", err)
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
