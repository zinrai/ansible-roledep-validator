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
	rolesPath := flag.String("roles", "roles", "Path to the Ansible roles directory")
	showVersion := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	if *playbookPath == "" || flag.NArg() > 0 {
		flag.Usage()
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
