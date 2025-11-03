package main

import (
	"flag"
	"fmt"
	"os"

	"mysql-tui-editor/server/internal/app"
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "./config/config.yml", "Path to configuration file")
	flag.Parse()

	// Print banner
	printBanner()

	// Create and run application
	application, err := app.New(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	// Ensure cleanup on exit
	defer func() {
		if err := application.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error during cleanup: %v\n", err)
		}
	}()

	// Run application
	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Application error: %v\n", err)
		os.Exit(1)
	}
}

func printBanner() {
	banner := `
╔═══════════════════════════════════════════════════╗
║                                                   ║
║        MySQL TUI Editor - Execution Server       ║
║                                                   ║
║        Educational SQL Learning Platform          ║
║                                                   ║
╚═══════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}
