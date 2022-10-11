package dormouse_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wabenet/dormouse"
)

var Tests = []struct {
	Name   string
	Args   []string
	Stdout string
	Stderr string
	Exit   int
}{
	{
		Name:   "greet steve",
		Args:   []string{"dormouse", "test/greet.yaml", "Steve", "--time", "morning"},
		Stdout: "Good morning, Steve!\n",
	},
	{
		Name:   "greet steve alt",
		Args:   []string{"dormouse", "test/greet.yaml", "Steve", "--time=morning"},
		Stdout: "Good morning, Steve!\n",
	},
	{
		Name:   "greet steve shorthand",
		Args:   []string{"dormouse", "test/greet.yaml", "Steve", "-t", "morning"},
		Stdout: "Good morning, Steve!\n",
	},
	{
		Name:   "greet steve informal",
		Args:   []string{"dormouse", "test/greet.yaml", "Steve", "--informal"},
		Stdout: "Hi Steve!\n",
	},
	{
		Name:   "greet steve informal shorthand",
		Args:   []string{"dormouse", "test/greet.yaml", "Steve", "-i"},
		Stdout: "Hi Steve!\n",
	},
	{
		Name:   "greet everyone",
		Args:   []string{"dormouse", "test/greet.yaml", "everyone"},
		Stdout: "Hello everyone!\n",
	},
	{
		Name:   "exit code",
		Args:   []string{"dormouse", "test/error.yaml"},
		Stdout: "This is an error\n",
		Exit:   127,
	},
}

func TestAll(t *testing.T) {
	t.Parallel()

	for _, tt := range Tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()

			var stdout strings.Builder
			var stderr strings.Builder

			d := &dormouse.Dormouse{
				Args:   tt.Args,
				Stdin:  bytes.NewReader(nil),
				Stdout: &stdout,
				Stderr: &stderr,
			}

			result := d.Execute()

			assert.Equal(t, tt.Exit, result)
			assert.Equal(t, tt.Stdout, stdout.String())
			assert.Equal(t, tt.Stderr, stderr.String())
		})
	}
}
