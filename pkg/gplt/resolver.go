package gplt

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Resolver interface {
	CanResolve(path string) bool
	Resolve(ctx *Context, path string) (string, error)
	DeriveName(path string) (string, error)
	CleanPath(ctx *Context, path string) (string, error)
	Trusted() bool
}

type ResolvedModule struct {
	Template	string
	Path		string
	Name		string
	Trusted		bool
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

func (m *ResolvedModule) GetHash() string {
	sum := sha256.Sum256([]byte(m.Template))
	hash := hex.EncodeToString(sum[:])

	return hash
}

func (c *ResolverChain) Resolve(ctx *Context, path string) (ResolvedModule, error) {
	for _, r := range c.resolvers {
		if r.CanResolve(path) {
			source, err := r.Resolve(ctx, path)

			if err != nil {
				return ResolvedModule{}, err
			}

			name, err := r.DeriveName(path)
			
			if err != nil {
				return ResolvedModule{}, err
			}

			clean, err := r.CleanPath(ctx, path)

			if err != nil {
				return ResolvedModule{}, err
			}

			return ResolvedModule{Template: source, Name: name, Path: clean, Trusted: r.Trusted()}, nil
		}
	}

	return ResolvedModule{}, fmt.Errorf("no resolver found for %q", path)
}