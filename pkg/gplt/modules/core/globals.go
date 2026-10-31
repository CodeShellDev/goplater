package core

import (
	"errors"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

type GlobalStore struct {
	data map[string]any
}

func NewGlobalStore() *GlobalStore {
	return &GlobalStore{
		data: map[string]any{},
	}
}

func (s *GlobalStore) Set(key string, value any) {
	s.data[key] = value.(string)
}

func (s *GlobalStore) Get(key string) (any, error) {
	value, exists := s.data[key]
	
	if !exists {
		return nil, errors.New("global '" + key + "' does not exist")
	}
	
	return value, nil
}

func (s *GlobalStore) Delete(key string) bool {
	delete(s.data, key)

	return true
}

func (s *GlobalStore) Keys() []string {
	keys := make([]string, 0, len(s.data))

	for k := range s.data {
		keys = append(keys, k)
	}

	return keys
}

func (s *GlobalStore) Has(key string) bool {
	_, exists := s.data[key]
	
	return exists
}

const globalStoreID = "globalStore"

// Defines a global value with a key.
//
// @param key string
// @param value any
// @returns error
var globalSetFunc = modules.NewFunc("globalSet", globalSet)

func globalSet(rt *gplt.Runtime, _ *gplt.Context, key string, value any) (string, error) {
	return "", SetGlobal(rt, key, value)
}

// Tries to retrieve a global value with a key.
//
// @param key string
// @returns any
// @returns error
var globalGetFunc = modules.NewFunc("globalGet", globalGet)

func globalGet(rt *gplt.Runtime, _ *gplt.Context, key string) (any, error)  {
	return GetGlobal(rt, key)
}

// Returns whether a given global variable exists.
//
// @param key string
// @returns bool
var globalHasFunc = modules.NewFunc("globalHas", globalHas)

func globalHas(rt *gplt.Runtime, _ *gplt.Context, key string) bool  {
	return HasGlobal(rt, key)
}

func SetGlobal(rt *gplt.Runtime, key string, value any) error {
	if !rt.HasStore(globalStoreID) {
		err := rt.RegisterStore(globalStoreID, NewGlobalStore())

		if err != nil {
			return errors.New("error registering global store: " + err.Error())
		}
	}

	s := rt.GetStore(globalStoreID)

	s.Set(key, value)

	return nil
}

func GetGlobal(rt *gplt.Runtime, key string) (any, error) {
	s := rt.GetStore(globalStoreID)

	return s.Get(key)
}

func HasGlobal(rt *gplt.Runtime, key string) bool {
	s := rt.GetStore(globalStoreID)

	return s.Has(key)
}