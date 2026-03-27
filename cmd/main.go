package main

import (
	"fmt"
	"os"

	"github.com/podutil/podutil/cmd/container"
	"github.com/podutil/podutil/internal/podman"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: podutil <command>")
	}

	switch args[1] {
	case "start":
		return handleStart(args[2:])
	case "help", "--help", "-h":
		return printHelp()
	default:
		return fmt.Errorf("unknown command: %s", args[1])
	}
}

func handleStart(args []string) error {
	if len(args) > 0 && args[0] != "" {
		id := args[0]
		client := podman.New()
		return client.StartContainer(id)
	}

	return container.StartCmd()
}

func printHelp() error {
	fmt.Println("podutil - CLI tool for Podman")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  podutil start [container_id]  Start a container")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  start    Start a container. If no ID provided, shows interactive list")
	return nil
}
