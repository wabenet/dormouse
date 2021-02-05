package dormouse

import (
	"fmt"
)

type Arguments []*Argument

type ArgumentValues map[string]string

type Argument struct {
	Name string `yaml:"name"`
}

func (as Arguments) Parse(args []string) (ArgumentValues, []string, error) {
	parsed := map[string]string{}

	for i, arg := range as {
		parsed[arg.Name] = args[i]
	}

	return parsed, args[len(as):], nil
}

func (as ArgumentValues) Get(name string) (string, error) {
	if value, ok := as[name]; ok {
		return value, nil
	}

	return "", fmt.Errorf("%w: undefined argument name: %s", ErrInvalidConfig, name)
}
