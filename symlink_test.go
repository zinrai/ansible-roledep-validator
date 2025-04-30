package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSymlinkHandling(t *testing.T) {
	// Create temporary directory structure for testing
	testDir, err := os.MkdirTemp("", "roles-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Create base directory for actual roles
	roleBaseDir := filepath.Join(testDir, "real-roles")
	if err := os.Mkdir(roleBaseDir, 0755); err != nil {
		t.Fatalf("Failed to create real roles dir: %v", err)
	}

	// 1. Create a valid role directory
	validRoleDir := filepath.Join(roleBaseDir, "valid-role")
	if err := os.Mkdir(validRoleDir, 0755); err != nil {
		t.Fatalf("Failed to create role dir: %v", err)
	}

	// 2. Create a file (not a directory)
	nonDirRolePath := filepath.Join(roleBaseDir, "non-dir-role")
	if f, err := os.Create(nonDirRolePath); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	} else {
		f.Close()
	}

	// Create directory for symlinks
	linkDir := filepath.Join(testDir, "roles")
	if err := os.Mkdir(linkDir, 0755); err != nil {
		t.Fatalf("Failed to create link dir: %v", err)
	}

	// Create three types of symlinks
	// a. Valid symlink to a directory
	if err := os.Symlink(validRoleDir, filepath.Join(linkDir, "valid-symlink")); err != nil {
		t.Fatalf("Failed to create valid symlink: %v", err)
	}

	// b. Broken symlink (points to non-existent path)
	if err := os.Symlink(filepath.Join(roleBaseDir, "nonexistent"), filepath.Join(linkDir, "broken-symlink")); err != nil {
		t.Fatalf("Failed to create broken symlink: %v", err)
	}

	// c. Symlink to a non-directory file
	if err := os.Symlink(nonDirRolePath, filepath.Join(linkDir, "non-dir-symlink")); err != nil {
		t.Fatalf("Failed to create non-dir symlink: %v", err)
	}

	// Run the test
	processor := NewRoleProcessor(false)
	roles, err := processor.getAllRoles(linkDir)

	// Verify the results
	if err != nil {
		t.Fatalf("Failed to get roles: %v", err)
	}

	// Verify that only valid symlinks are recognized as roles
	expectedRoles := []string{"valid-symlink"}
	if len(roles) != len(expectedRoles) {
		t.Errorf("Expected %d roles, got %d", len(expectedRoles), len(roles))
	}

	// Check if specific roles are included
	foundValidSymlink := false
	for _, role := range roles {
		if role == "valid-symlink" {
			foundValidSymlink = true
		}
	}

	if !foundValidSymlink {
		t.Errorf("Valid symlink role not found in results")
	}

	// Directly test the isValidSymlinkToDir function
	validDir, _ := processor.isValidSymlinkToDir(filepath.Join(linkDir, "valid-symlink"))
	if !validDir {
		t.Errorf("isValidSymlinkToDir failed to identify valid symlink to directory")
	}

	validDir, _ = processor.isValidSymlinkToDir(filepath.Join(linkDir, "broken-symlink"))
	if validDir {
		t.Errorf("isValidSymlinkToDir incorrectly identified broken symlink as valid")
	}

	validDir, _ = processor.isValidSymlinkToDir(filepath.Join(linkDir, "non-dir-symlink"))
	if validDir {
		t.Errorf("isValidSymlinkToDir incorrectly identified symlink to non-directory as valid")
	}
}
