package main

import (
	"fmt"
	"os"
	"path/filepath"
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

// Get all available roles from the roles directory, including symlinks to directories
func getAllRoles(rolesPath string) ([]string, error) {
	entries, err := os.ReadDir(rolesPath)
	if err != nil {
		return nil, fmt.Errorf("error reading roles directory: %v", err)
	}

	var roles []string
	for _, entry := range entries {
		if isDir(filepath.Join(rolesPath, entry.Name())) {
			roles = append(roles, entry.Name())
		}
	}
	return roles, nil
}

// Not DirEntry.IsDir: it reports false for a symlink to a directory, which
// would report a symlinked role as missing.
func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
