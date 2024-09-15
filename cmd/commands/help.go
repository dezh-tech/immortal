package commands

import "fmt"

const help = `Available commands:
Commands:
  run <path_to_config>      starts the application with the specified configuration file.
  help                      displays this help information.
  version                   displays the version of software.

Usage:
  <command> [options]

Examples:
  run config.yaml        run the application using config.yaml.
  help                   display this help message.
`

func HandleHelp(_ []string) {
	fmt.Print(help) //nolint
}
