package state

var (
	current = State{
		Archive: []string{},
		Audit: Audit{
			File:  []string{},
			Mount: []string{},
		},
		Lifeline: []Lifeline{},
		Tail:     []string{},
	}
)

type (
	State struct {
		Archive  []string
		Audit    Audit
		Lifeline []Lifeline
		Tail     []string
	}

	Audit struct {
		File  []string
		Mount []string
	}

	Lifeline struct {
		Name   string
		Script string
		Cycle  int
	}
)
