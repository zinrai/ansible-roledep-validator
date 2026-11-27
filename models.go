package main

import (
	"fmt"
	"strings"
)

// Accepts the two forms a role may be named in: a bare "rolename" and a
// "{ role: rolename }" mapping. Every other form is an error.
type RoleRef struct {
	Name string
}

func (r *RoleRef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var entry interface{}
	if err := unmarshal(&entry); err != nil {
		return err
	}

	name, ok := entry.(string)
	if mapping, isMapping := entry.(map[string]interface{}); isMapping {
		name, ok = mapping["role"].(string)
	}
	if !ok {
		return fmt.Errorf("role entry %v is neither a role name nor a mapping with a role key", entry)
	}
	r.Name = name
	return nil
}

// Represents a single play in an Ansible playbook
type Play struct {
	Roles []RoleRef `yaml:"roles"`
}

// Represents the structure of a role's meta/main.yml file
type MetaMain struct {
	Dependencies []RoleRef `yaml:"dependencies"`
}

// Extract the role names from a list of role references
func roleNames(refs []RoleRef) ([]string, error) {
	names := []string{}
	for _, r := range refs {
		if !isRoleDirName(r.Name) {
			return nil, fmt.Errorf("role %q is not the name of a directory in the roles directory", r.Name)
		}
		names = append(names, r.Name)
	}
	return names, nil
}

// Reports whether name can only refer to a directory directly under the roles
// directory. A path or a template resolves somewhere else.
func isRoleDirName(name string) bool {
	return name != "" && !strings.Contains(name, "/") && !strings.Contains(name, "{{")
}
