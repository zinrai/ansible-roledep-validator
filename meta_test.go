package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestDependencyResolution(t *testing.T) {
	testDir, err := os.MkdirTemp("", "roles-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	rolesDir := filepath.Join(testDir, "roles")
	if err := os.Mkdir(rolesDir, 0755); err != nil {
		t.Fatalf("Failed to create roles dir: %v", err)
	}

	// Create dependency chain: roleA -> roleB -> roleC, roleD
	roles := []string{"roleA", "roleB", "roleC", "roleD", "roleE"}

	for _, role := range roles {
		roleDir := filepath.Join(rolesDir, role)
		metaDir := filepath.Join(roleDir, "meta")

		if err := os.MkdirAll(metaDir, 0755); err != nil {
			t.Fatalf("Failed to create directory for %s: %v", role, err)
		}
	}

	// roleA depends on roleB
	writeMetaFile(t, filepath.Join(rolesDir, "roleA", "meta", "main.yml"), []string{"roleB"})

	// roleB depends on roleC and roleD
	writeMetaFile(t, filepath.Join(rolesDir, "roleB", "meta", "main.yml"), []string{"roleC", "roleD"})

	// roleD depends on roleE (creating a chain of dependencies)
	writeMetaFile(t, filepath.Join(rolesDir, "roleD", "meta", "main.yml"), []string{"roleE"})

	// roleC depends on roleA (creating a circular dependency)
	writeMetaFile(t, filepath.Join(rolesDir, "roleC", "meta", "main.yml"), []string{"roleA"})

	testCases := []struct {
		name         string
		initialRoles []string
		expectedDeps []string
	}{
		{
			name:         "Single role dependencies",
			initialRoles: []string{"roleA"},
			expectedDeps: []string{"roleA", "roleB", "roleC", "roleD", "roleE"},
		},
		{
			name:         "Multiple initial roles",
			initialRoles: []string{"roleB", "roleE"},
			expectedDeps: []string{"roleA", "roleB", "roleC", "roleD", "roleE"},
		},
		{
			name:         "With circular dependency",
			initialRoles: []string{"roleC"},
			expectedDeps: []string{"roleA", "roleB", "roleC", "roleD", "roleE"},
		},
		{
			name:         "Leaf dependency role",
			initialRoles: []string{"roleE"},
			expectedDeps: []string{"roleE"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			allDeps, err := getAllDependencies(tc.initialRoles, rolesDir)

			if err != nil {
				t.Fatalf("Failed to get dependencies: %v", err)
			}

			var actualDeps []string
			for role := range allDeps {
				actualDeps = append(actualDeps, role)
			}
			sort.Strings(actualDeps)
			sort.Strings(tc.expectedDeps)

			if !reflect.DeepEqual(actualDeps, tc.expectedDeps) {
				t.Errorf("Expected dependencies %v, got %v", tc.expectedDeps, actualDeps)
			}
		})
	}

}

func writeMetaFile(t *testing.T, path string, dependencies []string) {
	var metaYAML bytes.Buffer

	metaYAML.WriteString("dependencies:\n")
	for _, dep := range dependencies {
		metaYAML.WriteString(fmt.Sprintf("  - role: %s\n", dep))
	}

	if err := os.WriteFile(path, metaYAML.Bytes(), 0644); err != nil {
		t.Fatalf("Failed to write meta fil: %v", err)
	}
}

func TestUnreadableMetaIsAnError(t *testing.T) {
	rolesDir := t.TempDir()
	metaDir := filepath.Join(rolesDir, "roleA", "meta")
	if err := os.MkdirAll(metaDir, 0755); err != nil {
		t.Fatalf("Failed to create meta dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(metaDir, "main.yml"), []byte("dependencies: roleB\n"), 0644); err != nil {
		t.Fatalf("Failed to write meta file: %v", err)
	}

	allDeps, err := getAllDependencies([]string{"roleA"}, rolesDir)
	if err == nil {
		t.Fatalf("Expected an error, got dependencies %v", allDeps)
	}
}
