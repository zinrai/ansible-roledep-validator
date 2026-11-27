package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestPlaybookRoleForms(t *testing.T) {
	playbook := filepath.Join(t.TempDir(), "playbook.yml")
	content := "- hosts: all\n  roles:\n    - nginx\n    - role: docker.io\n      tags: container\n- hosts: db\n  roles:\n    - postgresql\n"
	if err := os.WriteFile(playbook, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write playbook: %v", err)
	}

	roles, err := getPlaybookRoles(playbook)
	if err != nil {
		t.Fatalf("Failed to get playbook roles: %v", err)
	}

	expectedRoles := []string{"nginx", "docker.io", "postgresql"}
	if !reflect.DeepEqual(roles, expectedRoles) {
		t.Errorf("Expected roles %v, got %v", expectedRoles, roles)
	}
}

func TestUnresolvableRoleEntryIsAnError(t *testing.T) {
	entries := map[string]string{
		"name key":    "- name: nginx",
		"no role key": "- tags: web",
		"path":        "- roles/nginx",
		"template":    `- "{{ webserver }}"`,
		"number":      "- 123",
		"null":        "- ~",
		"empty name":  `- role: ""`,
		"non-string":  "- role: 123",
	}

	for name, entry := range entries {
		t.Run(name+" in playbook", func(t *testing.T) {
			playbook := filepath.Join(t.TempDir(), "playbook.yml")
			content := "- hosts: all\n  roles:\n    " + entry + "\n"
			if err := os.WriteFile(playbook, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to write playbook: %v", err)
			}

			roles, err := getPlaybookRoles(playbook)
			if err == nil {
				t.Fatalf("Expected an error, got roles %v", roles)
			}
		})

		t.Run(name+" in meta", func(t *testing.T) {
			rolesDir := t.TempDir()
			metaDir := filepath.Join(rolesDir, "roleA", "meta")
			if err := os.MkdirAll(metaDir, 0755); err != nil {
				t.Fatalf("Failed to create meta dir: %v", err)
			}
			content := "dependencies:\n  " + entry + "\n"
			if err := os.WriteFile(filepath.Join(metaDir, "main.yml"), []byte(content), 0644); err != nil {
				t.Fatalf("Failed to write meta file: %v", err)
			}

			allDeps, err := getAllDependencies([]string{"roleA"}, rolesDir)
			if err == nil {
				t.Fatalf("Expected an error, got dependencies %v", allDeps)
			}
		})
	}
}
