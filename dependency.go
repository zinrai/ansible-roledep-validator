package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

// Extract role dependencies from meta/main.yml
func (rp *RoleProcessor) getDependencies(roleName string, rolesPath string) ([]string, error) {
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

	deps, err := roleNames(meta.Dependencies)
	if err != nil {
		return nil, err
	}

	if len(deps) > 0 {
		rp.Log("Found dependencies for role %s: %v\n", roleName, deps)
	}

	return deps, nil
}

// Recursively get all dependencies for a set of roles
func (rp *RoleProcessor) getAllDependencies(roles []string, rolesPath string) (map[string]bool, error) {
	rp.Log("Finding all dependencies for roles: %v\n", roles)

	allDeps := make(map[string]bool)
	for _, role := range roles {
		allDeps[role] = true
	}

	toProcess := append([]string{}, roles...)

	for len(toProcess) > 0 {
		role := toProcess[0]
		toProcess = toProcess[1:]

		deps, err := rp.getDependencies(role, rolesPath)
		if err != nil {
			return nil, fmt.Errorf("role %s: %v", role, err)
		}

		for _, dep := range deps {
			if !allDeps[dep] {
				rp.Log("Adding dependency: %s\n", dep)
				allDeps[dep] = true
				toProcess = append(toProcess, dep)
			}
		}
	}

	rp.Log("All required roles (including dependencies): %v\n", getMapKeys(allDeps))
	return allDeps, nil
}

// Get roles from a playbook
func (rp *RoleProcessor) getPlaybookRoles(playbookPath string) ([]string, error) {
	// Get absolute path
	absPath, err := filepath.Abs(playbookPath)
	if err != nil {
		return nil, fmt.Errorf("error resolving absolute path: %v", err)
	}
	rp.Log("Reading playbook: %s\n", absPath)

	// Check if file exists
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("error accessing playbook file: %v", err)
	}
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("playbook path is a directory, not a file")
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("error opening playbook file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading playbook file: %v", err)
	}

	var plays []Play
	err = yaml.Unmarshal(data, &plays)
	if err != nil {
		return nil, fmt.Errorf("error parsing playbook YAML: %v", err)
	}

	rp.Log("Found %d plays in playbook\n", len(plays))

	var roles []string
	for i, play := range plays {
		rp.Log("Processing play %d\n", i+1)

		names, err := roleNames(play.Roles)
		if err != nil {
			return nil, fmt.Errorf("play %d: %v", i+1, err)
		}
		roles = append(roles, names...)
	}

	rp.Log("Extracted %d roles from playbook: %v\n", len(roles), roles)
	return roles, nil
}
