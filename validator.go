package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func findMissingRoles(playbookPath, rolesPath string) ([]string, error) {
	playbookRoles, err := getPlaybookRoles(playbookPath)
	if err != nil {
		return nil, fmt.Errorf("error reading playbook: %v", err)
	}

	allAvailableRoles, err := getAllRoles(rolesPath)
	if err != nil {
		return nil, fmt.Errorf("error getting all roles: %v", err)
	}

	availableRolesMap := make(map[string]bool)
	for _, role := range allAvailableRoles {
		availableRolesMap[role] = true
	}

	allRequiredRoles, err := getAllDependencies(playbookRoles, rolesPath)
	if err != nil {
		return nil, fmt.Errorf("error getting dependencies: %v", err)
	}

	var missingRoles []string
	for role := range allRequiredRoles {
		if !availableRolesMap[role] {
			missingRoles = append(missingRoles, role)
		}
	}

	sort.Strings(missingRoles)
	return missingRoles, nil
}

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
