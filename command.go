package dormouse

import (
	"github.com/spf13/cobra"
)

type Command struct {
	Executable `yaml:",inline"`

	Description string              `yaml:"description"`
	Options     Options             `yaml:"options"`
	Arguments   Arguments           `yaml:"arguments"`
	Subcommands map[string]*Command `yaml:"subcommands"`
}

func (c *Command) ToCobraCommand(name string, r *result) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   name,
		Short: c.Description,
		Long:  c.Description,
		Args:  cobra.MinimumNArgs(len(c.Arguments)),
		Run: func(_ *cobra.Command, args []string) {
			if err := runCmd(c, args); err != nil {
				r.handleError(err)
			}
		},
	}

	for _, opt := range c.Options {
		if err := opt.Register(cmd); err != nil {
			return nil, err
		}
	}

	for n, sub := range c.Subcommands {
		subCmd, err := sub.ToCobraCommand(n, r)
		if err != nil {
			return nil, err
		}

		cmd.AddCommand(subCmd)
	}

	return cmd, nil
}

func runCmd(c *Command, args []string) error {
	templateArgs, remainder, err := c.Arguments.Parse(args)
	if err != nil {
		return err
	}

	templateOpts, err := c.Options.Parse()
	if err != nil {
		return err
	}

	cmd, err := c.Executable.Parse(templateOpts, templateArgs)
	if err != nil {
		return err
	}

	if err := cmd.Run(remainder); err != nil {
		return err
	}

	return nil
}
