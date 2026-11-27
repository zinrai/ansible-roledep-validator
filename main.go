package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	exitMissingRoles = 1
	exitFailure      = 2
)

func main() {
	playbookPath := flag.String("playbook", "", "Path to the Ansible playbook YAML file")
	rolesPath := flag.String("roles", "roles", "Path to the Ansible roles directory (default: roles)")
	showVersion := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	// If no playbook is provided, check if a positional argument is given
	args := flag.Args()
	if *playbookPath == "" && len(args) > 0 {
		*playbookPath = args[0]
	}

	if *playbookPath == "" {
		fmt.Fprintln(os.Stderr, "Usage: ansible-roledep-validator [-playbook PLAYBOOK_PATH] [-roles ROLES_PATH]")
		fmt.Fprintln(os.Stderr, "   or: ansible-roledep-validator PLAYBOOK_PATH")
		os.Exit(exitFailure)
	}

	missingRoles, err := findMissingRoles(*playbookPath, *rolesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(exitFailure)
	}

	for _, role := range missingRoles {
		fmt.Println(role)
	}

	if len(missingRoles) > 0 {
		os.Exit(exitMissingRoles)
	}
}
