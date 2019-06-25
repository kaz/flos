package state

import (
	"github.com/mattn/go-jsonpointer"
)

type (
	State struct {
		element interface{}
	}
)

func RootState() *State {
	return &State{store}
}

func (s *State) Get(path string) (*State, error) {
	elm, err := jsonpointer.Get(s.element, path)
	if err != nil {
		return nil, err
	}
	return &State{elm}, nil
}

func (s *State) List() []*State {
	ch := []*State{}
	for _, elm := range s.element.([]interface{}) {
		ch = append(ch, &State{elm})
	}
	return ch
}

func (s *State) Value() interface{} {
	return s.element
}
