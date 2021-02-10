package dormouse

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

const (
	Name        = "dormouse"
	Description = `
Dormouse is a stupidly simple tool, that builds and runs a simple CLI that wraps
existing tools and scripts based on a YAML configuration file.`

	ErrInvalidConfig Error = "invalid configuration"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type result struct {
	ExitCode int
	Message  string
}

func Execute(version string) int {
	r := &result{}

	rootCmd := &cobra.Command{
		Use:     fmt.Sprintf("%s file [args...]", Name),
		Short:   Description,
		Long:    Description,
		Version: version,
		Args:    cobra.MinimumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			if err := runRootCmd(r, args); err != nil {
				r.handleError(err)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		r.handleError(err)
	}

	if r.ExitCode != 0 {
		fmt.Fprintf(os.Stderr, "%s\n", r.Message)
	}

	return r.ExitCode
}

func runRootCmd(r *result, args []string) error {
	config, err := ReadConfigFromFile(args[0])
	if err != nil {
		return err
	}

	cmd, err := config.ToCobraCommand(fmt.Sprintf("%s %s", Name, args[0]), r)
	if err != nil {
		return err
	}

	cmd.SetArgs(args[1:])

	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}

	return nil
}

func (r *result) handleError(err error) {
	if err == nil {
		return
	}

	var e *exec.ExitError
	if errors.As(err, &e) {
		r.ExitCode = e.ExitCode()
	} else {
		r.ExitCode = 1
	}

	r.Message = err.Error()
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
