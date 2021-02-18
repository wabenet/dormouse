package dormouse

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

const (
	Description = `
Dormouse is a stupidly simple tool, that builds and runs a simple CLI that wraps
existing tools and scripts based on a YAML configuration file.`

	ErrInvalidConfig Error = "invalid configuration"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type Dormouse struct {
	Version string
	Args    []string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	exitCode int
	errorMsg string
}

func New(version string) *Dormouse {
	return &Dormouse{
		Version: version,
		Args:    os.Args,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
}

func (d *Dormouse) Execute() int {
	rootCmd := &cobra.Command{
		Use:     fmt.Sprintf("%s file [args...]", d.Args[0]),
		Short:   Description,
		Long:    Description,
		Version: d.Version,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config, err := ReadConfigFromFile(args[0])
			if err != nil {
				return err
			}

			cmd, err := config.ToCobraCommand(d, fmt.Sprintf("%s %s", d.Args[0], args[0]))
			if err != nil {
				return err
			}

			cmd.SetArgs(args[1:])

			if err := cmd.Execute(); err != nil {
				return fmt.Errorf("error executing command: %w", err)
			}

			return nil
		},
	}

	rootCmd.SetArgs(d.Args[1:])

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(d.Stderr, "%s\n", d.errorMsg)

		return 1
	}

	if d.exitCode != 0 {
		fmt.Fprintf(d.Stderr, "%s\n", d.errorMsg)
	}

	return d.exitCode
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
		d.errorMsg = err.Error()
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
