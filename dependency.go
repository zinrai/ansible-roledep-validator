package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
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

	deps := make([]string, len(meta.Dependencies))
	for i, dep := range meta.Dependencies {
		deps[i] = dep.Role
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
			rp.Log("Warning: error getting dependencies for role %s: %v\n", role, err)
			continue
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

// Extract a role from playbook data
func (rp *RoleProcessor) extractRoleFromPlaybook(roleValue interface{}, rolesList *[]string) {
	switch role := roleValue.(type) {
	case string:
		*rolesList = append(*rolesList, role)
		rp.Log("Found role (string): %s\n", role)
	case map[interface{}]interface{}:
		if roleName, ok := role["role"].(string); ok {
			*rolesList = append(*rolesList, roleName)
			rp.Log("Found role (map): %s\n", roleName)
		}
	default:
		rp.Log("Warning: unknown role format: %T\n", roleValue)
	}
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

		if play.Roles == nil {
			continue
		}

		rolesArray, ok := play.Roles.([]interface{})
		if !ok {
			rp.Log("Warning: roles section has unexpected format in play %d\n", i+1)
			continue
		}

		for _, item := range rolesArray {
			rp.extractRoleFromPlaybook(item, &roles)
		}
	}

	rp.Log("Extracted %d roles from playbook: %v\n", len(roles), roles)
	return roles, nil
}
