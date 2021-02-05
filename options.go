package dormouse

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Options []*Option

type OptionValues map[string]string

type Option struct {
	Name        string `yaml:"name"`
	Short       string `yaml:"short"`
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
	Required    bool   `yaml:"required"`

	value string
}

func (os Options) Parse() (OptionValues, error) {
	parsed := map[string]string{}

	for _, opt := range os {
		parsed[opt.Name] = opt.value
	}

	return parsed, nil
}

func (o *Option) Register(cmd *cobra.Command) error {
	flags := cmd.Flags()
	flags.VarP(o, o.Name, o.Short, o.Description)

	o.value = o.Default

	if o.Required {
		if err := cmd.MarkFlagRequired(o.Name); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}
	}

	return nil
}

func (o *Option) String() string {
	return o.value
}

func (o *Option) Set(value string) error {
	o.value = value

	return nil
}

func (o *Option) Type() string {
	return "string"
}

func (os OptionValues) Get(name string) (string, error) {
	if value, ok := os[name]; ok {
		return value, nil
	}

	return "", fmt.Errorf("%w: undefined option name: %s", ErrInvalidConfig, name)
}
