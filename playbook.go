package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

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

type Play struct {
	Roles []RoleRef `yaml:"roles"`
}

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

func isRoleDirName(name string) bool {
	return name != "" && !strings.Contains(name, "/") && !strings.Contains(name, "{{")
}
