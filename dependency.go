package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

// Extract role dependencies from meta/main.yml
func getDependencies(roleName string, rolesPath string) ([]string, error) {
	metaFile := filepath.Join(rolesPath, roleName, "meta", "main.yml")
	data, err := os.ReadFile(metaFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var meta MetaMain
	err = yaml.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}

	return roleNames(meta.Dependencies)
}

// Recursively get all dependencies for a set of roles
func getAllDependencies(roles []string, rolesPath string) (map[string]bool, error) {
	allDeps := make(map[string]bool)
	for _, role := range roles {
		allDeps[role] = true
	}

	toProcess := append([]string{}, roles...)

	for len(toProcess) > 0 {
		role := toProcess[0]
		toProcess = toProcess[1:]

		deps, err := getDependencies(role, rolesPath)
		if err != nil {
			return nil, fmt.Errorf("role %s: %v", role, err)
		}

		for _, dep := range deps {
			if !allDeps[dep] {
				allDeps[dep] = true
				toProcess = append(toProcess, dep)
			}
		}
	}

	return allDeps, nil
}

// Get roles from a playbook
func getPlaybookRoles(playbookPath string) ([]string, error) {
	data, err := os.ReadFile(playbookPath)
	if err != nil {
		return nil, err
	}

	var plays []Play
	err = yaml.Unmarshal(data, &plays)
	if err != nil {
		return nil, fmt.Errorf("error parsing playbook YAML: %v", err)
	}

	var roles []string
	for i, play := range plays {
		names, err := roleNames(play.Roles)
		if err != nil {
			return nil, fmt.Errorf("play %d: %v", i+1, err)
		}
		roles = append(roles, names...)
	}

	return roles, nil
}
