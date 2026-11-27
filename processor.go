package main

import (
	"fmt"
	"sort"
)

// Find roles that are required by the playbook but missing from the roles directory
func findMissingRoles(playbookPath, rolesPath string) ([]string, error) {
	// Step 1: Extract roles from playbook
	playbookRoles, err := getPlaybookRoles(playbookPath)
	if err != nil {
		return nil, fmt.Errorf("error reading playbook: %v", err)
	}

	// Step 2: Get all available roles
	allAvailableRoles, err := getAllRoles(rolesPath)
	if err != nil {
		return nil, fmt.Errorf("error getting all roles: %v", err)
	}

	// Convert available roles to a map for faster lookup
	availableRolesMap := make(map[string]bool)
	for _, role := range allAvailableRoles {
		availableRolesMap[role] = true
	}

	// Step 3: Get all dependencies
	allRequiredRoles, err := getAllDependencies(playbookRoles, rolesPath)
	if err != nil {
		return nil, fmt.Errorf("error getting dependencies: %v", err)
	}

	// Step 4: Find missing roles
	var missingRoles []string
	for role := range allRequiredRoles {
		if !availableRolesMap[role] {
			missingRoles = append(missingRoles, role)
		}
	}

	sort.Strings(missingRoles)
	return missingRoles, nil
}
