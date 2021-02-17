package dormouse_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/oclaussen/dormouse"
	"github.com/stretchr/testify/assert"
)

func FakeContext(args ...string) (*dormouse.Dormouse, *strings.Builder) {
	var stdout strings.Builder

	return &dormouse.Dormouse{
		Version: "test",
		Args:    args,
		Stdin:   bytes.NewReader(nil),
		Stdout:  &stdout,
		Stderr:  os.Stderr,
	}, &stdout
}

func TestGreetSteve(t *testing.T) {
	t.Parallel()

	d, stdout := FakeContext("dormouse", "test/greet.yaml", "--", "Steve", "--time", "morning")
	result := d.Execute()

	assert.Equal(t, 0, result)
	assert.Equal(t, "Good morning, Steve!\n", stdout.String())
}

func TestGreetEveryone(t *testing.T) {
	t.Parallel()

	d, stdout := FakeContext("dormouse", "test/greet.yaml", "--", "everyone")
	result := d.Execute()

	assert.Equal(t, 0, result)
	assert.Equal(t, "Hello everyone!\n", stdout.String())
}
