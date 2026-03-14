package core

import (
	"github.com/codeshelldev/goplater/pkg/templating"
	"github.com/codeshelldev/goplater/pkg/templating/modules"
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

func (s *GlobalStore) Get(key string) any {
	return s.data[key]
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

var globalSetFunc = modules.NewFunc("globalSet", globalSet)

func globalSet(rt *templating.Runtime, _ templating.Context, key string, value any) any  {
	SetGlobal(rt, key, value)
	return nil
}

var globalGetFunc = modules.NewFunc("globalGet", globalGet)

func globalGet(rt *templating.Runtime, _ templating.Context, key string) any  {
	return GetGlobal(rt, key)
}

func SetGlobal(rt *templating.Runtime, key string, value any) {
	if !rt.HasStore(globalStoreID) {
		err := rt.RegisterStore(globalStoreID, NewGlobalStore())

		if err != nil {
			panic("error registering global store: " + err.Error())
		}
	}

	s := rt.GetStore(globalStoreID)

	s.Set(key, value)
}

func GetGlobal(rt *templating.Runtime, key string) any {
	s := rt.GetStore(globalStoreID)

	return s.Get(key)
}