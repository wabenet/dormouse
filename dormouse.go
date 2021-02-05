package dormouse

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	yaml "gopkg.in/yaml.v2"
)

const (
	ErrInvalidArgument Error = "invalid argument"
	ErrInvalidConfig   Error = "invalid configuration"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

func Execute() int {
	if err := Run(); err != nil {
		var e *exec.ExitError
		if errors.As(err, &e) {
			return e.ExitCode()
		}

		fmt.Fprintf(os.Stderr, "%s\n", err)

		return 1
	}

	return 0
}

func Run() error {
	if len(os.Args) <= 1 {
		return fmt.Errorf("%w: first argument must be path to config file", ErrInvalidArgument)
	}

	config, err := ReadConfigFromFile(os.Args[1])
	if err != nil {
		return err
	}

	cmd, err := config.ToCobraCommand(path.Base(os.Args[0]))
	if err != nil {
		return err
	}

	cmd.SetArgs(os.Args[2:])

	return cmd.Execute()
}

func ReadConfigFromFile(path string) (*Command, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file '%s': %w", path, err)
	}

	return ReadConfig(bytes)
}

func ReadConfig(bytes []byte) (*Command, error) {
	var config Command
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("invalid yaml: %w", err)
	}

	return &config, nil
}
