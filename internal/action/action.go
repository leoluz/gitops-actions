package action

import (
	"log"
	"os"
	"strings"

	"github.com/leoluz/gitops-actions/internal/git"
)

type Action interface {
	Execute() error
}

type ActionCreator interface {
	NewAction(file *git.File) Action
}

type Registry struct {
	db map[string]ActionCreator
}

type RegistryParams struct {
	TwitterConsumerKey       string
	TwitterConsumerSecret    string
	TwitterAccessToken       string
	TwitterAccessTokenSecret string
}

func NewRegistry() *Registry {
	return &Registry{
		db: make(map[string]ActionCreator),
	}
}

func (r *Registry) Add(name string, ac ActionCreator) {
	r.db[name] = ac
}

func (r *Registry) Get(name string) ActionCreator {
	if ac, ok := r.db[name]; ok {
		return ac
	}
	return nil
}

func BuildActions(reg *Registry, actionsDir string, files []*git.File) ([]Action, error) {
	log.Printf("building actions for %d files", len(files))
	actions := []Action{}
	for _, file := range files {
		if strings.HasPrefix(file.GetName(), actionsDir) {
			dirs := strings.Split(file.GetName(), string(os.PathSeparator))
			if len(dirs) > 2 {
				actionName := dirs[1]
				actionCreator := reg.Get(actionName)
				if actionCreator == nil {
					log.Printf("action %q not found in registry: skipping file %s", actionName, file.GetName())
					continue
				}
				action := actionCreator.NewAction(file)
				actions = append(actions, action)
				log.Printf("%s action created for %s file %s", actionName, file.GetStatus().String(), file.GetName())
			}
		}
	}
	return actions, nil
}
