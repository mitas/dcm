package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/mitas/dcm/internal/model"
)

// Config holds all application configuration
type Config struct {
	// RootPath is the directory to search for docker-compose projects
	RootPath string
	// Action is the action to perform
	Action model.ActionType
	// TargetProject is the name of the project to act on (if specified)
	TargetProject string
	// ActionAll indicates if action should be performed on all projects
	ActionAll bool
}

// ParseConfig parses command line flags and returns the application config
func ParseConfig() (*Config, error) {
	config := &Config{}

	// Define flags
	rootPath := flag.String("path", "", "Root path to search for docker-compose projects")
	startAll := flag.Bool("start-all", false, "Start all docker-compose projects")
	stopAll := flag.Bool("stop-all", false, "Stop all docker-compose projects")
	start := flag.String("start", "", "Start specific docker-compose project (provide project name)")
	stop := flag.String("stop", "", "Stop specific docker-compose project (provide project name)")
	status := flag.String("status", "", "Check status of specific docker-compose project (provide project name)")
	statusAll := flag.Bool("status-all", false, "Check status of all docker-compose projects")
	list := flag.Bool("list", false, "List all docker-compose projects")

	flag.Parse()

	// Validate root path
	if *rootPath == "" {
		return nil, fmt.Errorf("root path is required, please provide using -path flag")
	}
	config.RootPath = *rootPath

	// Determine action type
	switch {
	case *list:
		config.Action = model.ActionList
	case *startAll:
		config.Action = model.ActionStart
		config.ActionAll = true
	case *stopAll:
		config.Action = model.ActionStop
		config.ActionAll = true
	case *statusAll:
		config.Action = model.ActionStatus
		config.ActionAll = true
	case *start != "":
		config.Action = model.ActionStart
		config.TargetProject = *start
	case *stop != "":
		config.Action = model.ActionStop
		config.TargetProject = *stop
	case *status != "":
		config.Action = model.ActionStatus
		config.TargetProject = *status
	default:
		return nil, fmt.Errorf("no action specified. Use -list, -start-all, -stop-all, -status-all, -start=PROJECT, -stop=PROJECT, or -status=PROJECT")
	}

	return config, nil
}

// PrintUsage prints usage information
func PrintUsage() {
	fmt.Println("Usage: dcm -path ROOT_PATH [flags]")
	flag.PrintDefaults()
	os.Exit(1)
}