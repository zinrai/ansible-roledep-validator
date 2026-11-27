package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

// Represents the structure of a role's meta/main.yml file
type MetaMain struct {
	Dependencies []RoleRef `yaml:"dependencies"`
}

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
	toProcess := append([]string{}, roles...)

	for len(toProcess) > 0 {
		role := toProcess[0]
		toProcess = toProcess[1:]
		if allDeps[role] {
			continue
		}
		allDeps[role] = true

		deps, err := getDependencies(role, rolesPath)
		if err != nil {
			return nil, fmt.Errorf("role %s: %v", role, err)
		}
		toProcess = append(toProcess, deps...)
	}

	return allDeps, nil
}
