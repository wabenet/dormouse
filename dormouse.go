package dormouse

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	yaml "gopkg.in/yaml.v2"
)

const (
	ErrInvalidConfig    Error = "invalid configuration"
	ErrInvalidArguments Error = "invalid arguments"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type Dormouse struct {
	Args []string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	exitCode int
}

func New(version string) *Dormouse {
	return &Dormouse{
		Args:     os.Args,
		Stdin:    os.Stdin,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		exitCode: 0,
	}
}

func (d *Dormouse) Execute() int {
	if len(d.Args) < 2 {
		return d.fail(fmt.Errorf("%w: first argument must be path to config file", ErrInvalidArguments))
	}

	cmd, err := ReadConfigFromFile(d.Args[1])
	if err != nil {
		return d.fail(err)
	}

	if err := cmd.Execute(d, d.Args[2:]); err != nil {
		return d.fail(err)
	}

	return d.exitCode
}

func (d *Dormouse) fail(err error) int {
	fmt.Fprintf(d.Stderr, "%s\n", err.Error())

	return 1
}

func (d *Dormouse) Exec(path string, args ...string) error {
	// #nosec G204 // Because that is the whole point if this tool
	cmd := exec.Command(path, args...)
	cmd.Stdin = d.Stdin
	cmd.Stdout = d.Stdout
	cmd.Stderr = d.Stderr

	err := cmd.Run()

	if err == nil {
		return nil
	}

	var e *exec.ExitError
	if errors.As(err, &e) {
		d.exitCode = e.ExitCode()

		return nil
	}

	return fmt.Errorf("could not execute: %w", err)
}

func ReadConfigFromFile(path string) (*Command, error) {
	bytes, err := os.ReadFile(path)
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
