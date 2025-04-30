package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Handles role processing with optional logging
type RoleProcessor struct {
	logOutput io.Writer
}

// Creates a new processor with the given log output
func NewRoleProcessor(verbose bool) *RoleProcessor {
	var logOutput io.Writer
	if verbose {
		logOutput = os.Stdout
	} else {
		logOutput, _ = os.Open(os.DevNull)
	}
	return &RoleProcessor{logOutput: logOutput}
}

// Prints a message if verbose mode is enabled
func (rp *RoleProcessor) Log(format string, args ...interface{}) {
	fmt.Fprintf(rp.logOutput, format, args...)
}

// Find roles that are required by the playbook but missing from the roles directory
func (rp *RoleProcessor) findMissingRoles(playbookPath, rolesPath string) ([]string, error) {
	// Step 1: Extract roles from playbook
	playbookRoles, err := rp.getPlaybookRoles(playbookPath)
	if err != nil {
		return nil, fmt.Errorf("error reading playbook: %v", err)
	}

	// Step 2: Get all available roles
	allAvailableRoles, err := rp.getAllRoles(rolesPath)
	if err != nil {
		return nil, fmt.Errorf("error getting all roles: %v", err)
	}

	// Convert available roles to a map for faster lookup
	availableRolesMap := make(map[string]bool)
	for _, role := range allAvailableRoles {
		availableRolesMap[role] = true
	}

	// Step 3: Get all dependencies
	allRequiredRoles, err := rp.getAllDependencies(playbookRoles, rolesPath)
	if err != nil {
		return nil, fmt.Errorf("error getting dependencies: %v", err)
	}

	// Step 4: Find missing roles
	var missingRoles []string
	for role := range allRequiredRoles {
		if !availableRolesMap[role] {
			rp.Log("Role '%s' is missing (not found in roles directory)\n", role)
			missingRoles = append(missingRoles, role)
		}
	}

	sort.Strings(missingRoles)
	return missingRoles, nil
}

// Helper function to get keys from a map
func getMapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
