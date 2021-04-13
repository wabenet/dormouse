package dormouse

import (
	"fmt"
	"strings"
)

const endOfOptions = "--"

type Arguments struct {
	Positionals Positionals `yaml:"arguments"`
	Options     Options     `yaml:"options"`
}

type Positionals []*Positional

type Positional struct {
	Name string `yaml:"name"`

	value string `yaml:"-"`
}

type Options []*Option

type Option struct {
	Name     string `yaml:"name"`
	Short    string `yaml:"short"`
	Default  string `yaml:"default"`
	Required bool   `yaml:"required"`

	value string `yaml:"-"`
}

type Values struct {
	args             []string
	positionalByName map[string]*Positional
	optionByName     map[string]*Option
	optionByShort    map[byte]*Option
}

func NewValues() *Values {
	return &Values{
		args:             []string{},
		positionalByName: map[string]*Positional{},
		optionByName:     map[string]*Option{},
		optionByShort:    map[byte]*Option{},
	}
}

func (vs *Values) Positional(name string) (string, error) {
	if arg, ok := vs.positionalByName[name]; ok {
		return arg.value, nil
	}

	return "", fmt.Errorf("%w: undefined argument name: %s", ErrInvalidConfig, name)
}

func (vs *Values) Option(name string) (string, error) {
	if opt, ok := vs.optionByName[name]; ok {
		return opt.value, nil
	}

	return "", fmt.Errorf("%w: undefined option name: %s", ErrInvalidConfig, name)
}

func (vs *Values) addOption(opt *Option) error {
	opt.value = opt.Default

	if _, exist := vs.optionByName[opt.Name]; exist {
		return fmt.Errorf("%w: flag redefined: %s", ErrInvalidConfig, opt.Name)
	}

	vs.optionByName[opt.Name] = opt

	if opt.Short == "" {
		return nil
	}

	if len(opt.Short) > 1 {
		return fmt.Errorf("%w: %q shorthand must be a single character", ErrInvalidConfig, opt.Short)
	}

	short := opt.Short[0]
	if other, exist := vs.optionByShort[short]; exist {
		return fmt.Errorf("%w: shorthand %q is already used for option %q", ErrInvalidConfig, short, other.Name)
	}

	vs.optionByShort[short] = opt

	return nil
}

func (as *Arguments) Parse(args []string) (*Values, error) {
	vs := NewValues()

	for _, opt := range as.Options {
		if err := vs.addOption(opt); err != nil {
			return nil, err
		}
	}

	if err := vs.parseArgs(args); err != nil {
		return nil, err
	}

	if len(vs.args) < len(as.Positionals) {
		return nil, fmt.Errorf("%w: requires at least %d arg(s)", ErrInvalidArguments, len(as.Positionals))
	}

	for i, arg := range as.Positionals {
		arg.value = vs.args[i]
		vs.positionalByName[arg.Name] = arg
	}

	vs.args = vs.args[len(as.Positionals):]

	for _, opt := range as.Options {
		if opt.Required && opt.value == "" {
			return nil, fmt.Errorf("%w: option %s is required", ErrInvalidArguments, opt.Name)
		}
	}

	return vs, nil
}

func (vs *Values) parseArgs(args []string) error {
	if len(args) == 0 {
		return nil
	}

	current, remainder := args[0], args[1:]

	if current == endOfOptions {
		vs.args = append(vs.args, remainder...)

		return nil
	}

	if ok, name, err := isLongOption(current); ok {
		if err != nil {
			return err
		}

		split := strings.SplitN(name, "=", 2)
		name = split[0]
		opt, exists := vs.optionByName[name]

		if !exists {
			vs.args = append(vs.args, current)

			return vs.parseArgs(remainder)
		}

		if len(split) == 2 { // --name=value
			opt.value = split[1]
		} else if len(remainder) > 0 { // --name value
			opt.value, remainder = remainder[0], remainder[1:]
		} else { // --name
			return fmt.Errorf("%w: value required for %s", ErrInvalidArguments, current)
		}

		return vs.parseArgs(remainder)
	}

	if ok, short, err := isShortOption(current); ok {
		if err != nil {
			return err
		}

		opt, exists := vs.optionByShort[short]
		if !exists {
			vs.args = append(vs.args, current)

			return vs.parseArgs(remainder)
		}

		if len(current) > 3 && current[2] == '=' { // -s=value
			opt.value = current[3:]
		} else if len(current) > 2 { // -svalue
			opt.value = current[2:]
		} else if len(remainder) > 0 { // -s value
			opt.value, remainder = remainder[0], remainder[1:]
		} else { // -s
			return fmt.Errorf("%w: value required for %q in %s", ErrInvalidArguments, short, current)
		}

		return vs.parseArgs(remainder)
	}

	vs.args = append(vs.args, current)

	return vs.parseArgs(remainder)
}

func isLongOption(arg string) (bool, string, error) {
	if len(arg) < 3 || !strings.HasPrefix(arg, "--") {
		return false, "", nil
	}

	if arg[3] == '-' || arg[3] == '=' {
		return true, "", fmt.Errorf("%w: invalid option syntax: %s", ErrInvalidArguments, arg)
	}

	return true, arg[2:], nil
}

func isShortOption(arg string) (bool, byte, error) {
	if len(arg) < 2 || arg[0] != '-' || arg[1] == '-' {
		return false, 0, nil
	}

	if arg[1] == '=' {
		return true, 0, fmt.Errorf("%w: invalid short option syntax: %s", ErrInvalidArguments, arg)
	}

	return true, arg[1], nil
}
