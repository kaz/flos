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
	src, ok := s.element.([]interface{})
	if !ok {
		logger.Println("target is not array! ", s)
		return []*State{}
	}

	result := []*State{}
	for _, elm := range src {
		result = append(result, &State{elm})
	}
	return result
}

func (s *State) Value() interface{} {
	return s.element
}
