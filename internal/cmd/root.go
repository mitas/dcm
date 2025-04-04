package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mitas/dcm/internal/manager"
	"github.com/mitas/dcm/pkg/formatter"
)

var (
	rootPath   string
	configPath string
)

// NewRootCmd creates the root command for the application
func NewRootCmd(projectManager *manager.Manager, outputFormatter *formatter.Formatter) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "dcm",
		Short: "Docker Compose Manager - Manage multiple docker-compose projects",
		Long: `Docker Compose Manager (dcm) is a tool for finding and managing 
multiple docker-compose projects in a directory structure.

It allows you to list, start, stop, and check the status of docker-compose 
projects in a given directory.`,
		SilenceUsage: true,
	}

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&rootPath, "path", "p", "", "Root path to search for docker-compose projects")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to config file (default is ~/.config/dcm/config.yaml)")

	// Add subcommands
	rootCmd.AddCommand(newListCmd(projectManager, outputFormatter))
	rootCmd.AddCommand(newStartCmd(projectManager, outputFormatter))
	rootCmd.AddCommand(newStopCmd(projectManager, outputFormatter))
	rootCmd.AddCommand(newStatusCmd(projectManager, outputFormatter))

	// Add managed project commands
	rootCmd.AddCommand(newListManagedCmd(projectManager, outputFormatter))
	rootCmd.AddCommand(newAddManagedCmd(projectManager, outputFormatter))
	rootCmd.AddCommand(newRemoveManagedCmd(projectManager, outputFormatter))

	return rootCmd
}
