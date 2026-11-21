package test

import (
	"fmt"

	"github.com/codeshelldev/goplater/pkg/gplt"
	"github.com/codeshelldev/goplater/pkg/gplt/modules"
)

const storeID = "testStore"

type Entry struct {
	Label string
	Value any
}

type Store struct {
	entries []Entry
	byLabel map[string]any
}

func NewStore() *Store {
	return &Store{byLabel: map[string]any{}}
}


func (s *Store) Set(label string, value any) {
	s.entries = append(s.entries, Entry{label, value})

	if label != "" {
		s.byLabel[label] = value
	}
}

func (s *Store) Get(label string) (any, error) {
	v, ok := s.byLabel[label]
	if !ok {
		return nil, fmt.Errorf("no emit with label %q", label)
	}

	return v, nil
}

func (s *Store) Has(label string) bool {
	_, ok := s.byLabel[label]
	return ok
}

func (s *Store) Delete(label string) bool {
	_, ok := s.byLabel[label]
	if !ok {
		return false
	}

	delete(s.byLabel, label)
	return true
}

func (s *Store) Keys() []string {
	keys := make([]string, 0, len(s.byLabel))

	for k := range s.byLabel {
		keys = append(keys, k)
	}

	return keys
}

func (s *Store) Entries() []Entry {
	return append([]Entry(nil), s.entries...)
}

func getStore(rt *gplt.Runtime) (*Store, error) {
	if !rt.HasStore(storeID) {
		err := rt.RegisterStore(storeID, NewStore())
		if err != nil {
			return nil, fmt.Errorf("error registering test store: %w", err)
		}
	}

	s, ok := rt.GetStore(storeID).(*Store)

	if !ok {
		return nil, fmt.Errorf("test store has wrong type")
	}

	return s, nil
}

func GetResults(rt *gplt.Runtime) (*Store, bool) {
	if !rt.HasStore(storeID) {
		return nil, false
	}

	s, ok := rt.GetStore(storeID).(*Store)
	return s, ok
}

var Module = modules.NewModule(emitFunc)

// Records a value, optionally with a label.
//
// Not bundled with goplater, used internally for unit and integration tests.
//
// @param label string?
// @param value any
// @returns error
var emitFunc = modules.NewFunc("emit", emitValue, emitLabeled)

func emitValue(rt *gplt.Runtime, _ *gplt.Context, value any) (string, error) {
	return emitLabeled(rt, nil, "", value)
}

func emitLabeled(rt *gplt.Runtime, _ *gplt.Context, label string, value any) (string, error) {
	s, err := getStore(rt)

	if err != nil {
		return "", err
	}

	s.Set(label, value)

	return "", nil
}