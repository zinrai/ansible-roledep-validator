package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	playbookPath := flag.String("playbook", "", "Path to the Ansible playbook YAML file")
	rolesPath := flag.String("roles", "roles", "Path to the Ansible roles directory (default: roles)")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	// If no playbook is provided, check if a positional argument is given
	args := flag.Args()
	if *playbookPath == "" && len(args) > 0 {
		*playbookPath = args[0]
	}

	if *playbookPath == "" {
		fmt.Fprintln(os.Stderr, "Usage: ansible-roledep-validator [-playbook PLAYBOOK_PATH] [-roles ROLES_PATH] [-verbose]")
		fmt.Fprintln(os.Stderr, "   or: ansible-roledep-validator PLAYBOOK_PATH")
		os.Exit(1)
	}

	// Create role processor with verbose setting
	rp := NewRoleProcessor(*verbose)

	// Display header information if verbose
	if *verbose {
		fmt.Printf("\n=== Ansible Role Dependency Validator ===\n")
		fmt.Printf("Playbook: %s\n", *playbookPath)
		fmt.Printf("Roles directory: %s\n\n", *rolesPath)

		cwd, _ := os.Getwd()
		fmt.Printf("Current working directory: %s\n\n", cwd)
	}

	// Find missing roles using the processor
	missingRoles, err := rp.findMissingRoles(*playbookPath, *rolesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Output results based on verbose setting
	if *verbose {
		fmt.Printf("\n=== VALIDATION RESULTS ===\n")
		if len(missingRoles) > 0 {
			fmt.Println("Missing roles:")
			for _, role := range missingRoles {
				fmt.Printf("- %s\n", role)
			}
		} else {
			fmt.Println("All roles are present.")
		}
	} else {
		// Just print the missing role names
		for _, role := range missingRoles {
			fmt.Println(role)
		}
	}

	// Exit with error code if missing roles found
	if len(missingRoles) > 0 {
		os.Exit(1)
	}
}
