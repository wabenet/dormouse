module github.com/oclaussen/dormouse

go 1.16

// Adds support to access the unknown flags from pflag.
// Necessary so we can pass these on to the target script.
replace github.com/spf13/pflag => github.com/ijc/spf13-pflag v0.0.0-20190125153223-1fb1288bf36d

require (
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.3.0
	gopkg.in/yaml.v2 v2.4.0
)
