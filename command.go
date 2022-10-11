package dormouse

type Commands map[string]*Command

type Command struct {
	Subcommands Commands `yaml:"subcommands"`

	Executable `yaml:",inline"`
	Arguments  `yaml:",inline"`
}

func (c *Command) Execute(d *Dormouse, args []string) error {
	if len(args) == 0 {
		return c.run(d, args)
	}

	if sub, ok := c.Subcommands[args[0]]; ok {
		return sub.Execute(d, args[1:])
	}

	return c.run(d, args)
}

func (c *Command) run(d *Dormouse, args []string) error {
	values, err := c.Arguments.Parse(args)
	if err != nil {
		return err
	}

	ex, err := c.Executable.Parse(values)
	if err != nil {
		return err
	}

	return d.Exec(ex, values.args)
}
