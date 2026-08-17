package templating

import "errors"

type Resolver interface {
	CanResolve(path string) bool
	Resolve(path string) (string, error)
	DeriveName(path string) (string, error)
	Trusted() bool
}

type ResolvedModule struct {
	Template  	string
	Name 		string
	Trusted 	bool
}

type ResolverChain struct {
	resolvers []Resolver
}

const IMPORT_FILE_EXTENSION = "gplt"

func NewResolverChain(resolvers ...Resolver) *ResolverChain {
	return &ResolverChain{resolvers: resolvers}
}

func (c *ResolverChain) Add(r Resolver) {
	c.resolvers = append(c.resolvers, r)
}

func (c *ResolverChain) Resolve(path string) (ResolvedModule, error) {
	for _, r := range c.resolvers {
		if r.CanResolve(path) {
			source, err := r.Resolve(path)

			if err != nil {
				return ResolvedModule{}, err
			}

			name, err := r.DeriveName(path)
			
			if err != nil {
				return ResolvedModule{}, err
			}

			return ResolvedModule{Template: source, Name: name, Trusted: r.Trusted()}, nil
		}
	}

	return ResolvedModule{}, errors.New("no resolver found for " + path)
}