package main

import (
	"os"

	"github.com/wabenet/dormouse"
)

var version = "latest"

func main() {
	os.Exit(dormouse.New(version).Execute())
}
