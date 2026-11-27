package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Get all available roles from the roles directory, including symlinks to directories
func getAllRoles(rolesPath string) ([]string, error) {
	entries, err := os.ReadDir(rolesPath)
	if err != nil {
		return nil, fmt.Errorf("error reading roles directory: %v", err)
	}

	var roles []string
	for _, entry := range entries {
		if isDir(filepath.Join(rolesPath, entry.Name())) {
			roles = append(roles, entry.Name())
		}
	}
	return roles, nil
}

// Not DirEntry.IsDir: it reports false for a symlink to a directory, which
// would report a symlinked role as missing.
func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
