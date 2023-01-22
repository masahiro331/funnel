package main

import (
	"os"

	"github.com/masahiro331/funnel/cmd/commander"
	"github.com/masahiro331/funnel/cmd/operator"
	"github.com/masahiro331/funnel/cmd/pod"
)

func main() {
	subcommand := os.Args[1]
	switch subcommand {
	case "commander":
		commander.Run(os.Args[2], os.Args[3])
	case "operator":
		operator.Run()
	case "pod":
		pod.Run()
	}
}
