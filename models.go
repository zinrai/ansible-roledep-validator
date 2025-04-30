package main

// Represents a single play in an Ansible playbook
type Play struct {
	Roles interface{} `yaml:"roles"`
}

// Represents the structure of a role's meta/main.yml file
type MetaMain struct {
	Dependencies []struct {
		Role string `yaml:"role"`
	} `yaml:"dependencies"`
}
