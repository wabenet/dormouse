package dormouse

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

type Executable struct {
	Exec   []string `yaml:"exec"`
	Script string   `yaml:"script"`
}

type ExecutableCommand struct {
	Path string
	Args []string
}

func (e *Executable) Parse(opts OptionValues, args ArgumentValues) (*ExecutableCommand, error) {
	if len(e.Exec) > 0 && len(e.Script) > 0 {
		return nil, fmt.Errorf("%w: only one of 'exec' and 'script' is allowed", ErrInvalidConfig)
	}

	fs := template.FuncMap{
		"option": opts.Get,
		"arg":    args.Get,
		"which":  exec.LookPath,
	}

	execArgs := []string{}

	if len(e.Exec) > 0 {
		for _, arg := range e.Exec {
			value, err := runTemplate(arg, fs)
			if err != nil {
				return nil, err
			}

			execArgs = append(execArgs, value)
		}
	}

	if len(e.Script) > 0 {
		value, err := runTemplate(e.Script, fs)
		if err != nil {
			return nil, err
		}

		execArgs = append(execArgs, "/bin/sh", "-c", value, "--")
	}

	return &ExecutableCommand{Path: execArgs[0], Args: execArgs[1:]}, nil
}

func runTemplate(text string, fs template.FuncMap) (string, error) {
	templ, err := template.New("script").Funcs(fs).Parse(text)
	if err != nil {
		return "", fmt.Errorf("could not parse template string: %w", err)
	}

	var buffer bytes.Buffer

	if err := templ.Execute(&buffer, nil); err != nil {
		return "", fmt.Errorf("could not evaluate template: %w", err)
	}

	return buffer.String(), nil
}

func (e *ExecutableCommand) Run(args []string) error {
	// #nosec G204 // Because that is the whole point if this tool
	cmd := exec.Command(e.Path, append(e.Args, args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
