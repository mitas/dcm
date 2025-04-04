package main

import (
	"fmt"
	"os"

	"github.com/mitas/dcm/internal/cmd"
	"github.com/mitas/dcm/internal/manager"
	"github.com/mitas/dcm/pkg/formatter"
)

func main() {
	// Initialize components
	cmdExecutor := &manager.DefaultCommandExecutor{}
	projectManager := manager.NewManager(cmdExecutor)
	outputFormatter := formatter.NewFormatter()

	// Initialize root command
	rootCmd := cmd.NewRootCmd(projectManager, outputFormatter)

	// Execute the application
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
