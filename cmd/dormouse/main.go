package main

import (
	"os"

	"github.com/oclaussen/dormouse"
)

var version = "latest"

func main() {
	os.Exit(dormouse.Execute(version))
}
