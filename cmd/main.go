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
		return printHelp()
	}

	switch args[1] {
	case "start":
		return handleStart(args[2:])
	case "stop":
		return handleStop(args[2:])
	case "ps":
		return handlePs()
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

func handleStop(args []string) error {
	if len(args) > 0 && args[0] != "" {
		id := args[0]
		client := podman.New()
		return client.StopContainer(id)
	}

	return container.StopCmd()
}

func handlePs() error {
	return container.PsCmd()
}

func printHelp() error {
	fmt.Println("podutil - CLI tool for Podman")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  podutil start [container_id]  Start a container")
	fmt.Println("  podutil stop [container_id]   Stop a container")
	fmt.Println("  podutil ps                    List running containers")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  start [id]  Start a container (with TUI if no ID)")
	fmt.Println("  stop [id]   Stop a container (with TUI if no ID)")
	fmt.Println("  ps          List running containers with details")
	fmt.Println("  help        Show this help message")
	return nil
}
