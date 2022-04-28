package action

type Action interface {
	Execute() error
}

type ActionCreator interface {
	NewAction(fileContent string) Action
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
