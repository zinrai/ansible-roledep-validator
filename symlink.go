package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Check if a path is a valid symlink that points to a directory
func (rp *RoleProcessor) isValidSymlinkToDir(path string) (bool, string) {
	// Get the symlink target
	target, err := os.Readlink(path)
	if err != nil {
		return false, "" // Not a symlink or error reading it
	}

	// If target is not absolute, make it absolute relative to the symlink's directory
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(path), target)
	}

	// Check if the target exists and is a directory
	targetInfo, err := os.Stat(target)
	if err != nil {
		return false, target // Target doesn't exist
	}

	return targetInfo.IsDir(), target
}

// Process a directory entry for role detection
func (rp *RoleProcessor) processDirectoryEntry(entry os.DirEntry, absPath string, roles *[]string) {
	entryPath := filepath.Join(absPath, entry.Name())
	fileInfo, err := os.Lstat(entryPath)
	if err != nil {
		rp.Log("Error reading entry %s: %v\n", entry.Name(), err)
		return
	}

	// Handle regular directory
	if fileInfo.IsDir() {
		*roles = append(*roles, entry.Name())
		rp.Log("Found role (directory): %s\n", entry.Name())
		return
	}

	// Handle symlink
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		isValidDir, targetPath := rp.isValidSymlinkToDir(entryPath)
		if isValidDir {
			*roles = append(*roles, entry.Name())
			rp.Log("Found role (symlink): %s -> %s\n", entry.Name(), targetPath)
		} else if targetPath != "" {
			rp.Log("Symlink %s points to %s which is not a directory\n", entry.Name(), targetPath)
		} else {
			rp.Log("Found broken symlink: %s\n", entry.Name())
		}
	}
}

// Print debug info about directory contents
func (rp *RoleProcessor) printDirectoryDebugInfo(entries []os.DirEntry, absPath string) {
	rp.Log("Directory contents:\n")
	for _, entry := range entries {
		entryPath := filepath.Join(absPath, entry.Name())
		fileInfo, err := os.Lstat(entryPath)
		if err != nil {
			rp.Log("- %s (error reading info)\n", entry.Name())
			continue
		}

		if fileInfo.Mode()&os.ModeSymlink != 0 {
			target, _ := os.Readlink(entryPath)
			rp.Log("- %s (symlink -> %s)\n", entry.Name(), target)
		} else {
			rp.Log("- %s (isDir: %t)\n", entry.Name(), fileInfo.IsDir())
		}
	}
}

// Get all available roles from the roles directory, properly handling symlinks
func (rp *RoleProcessor) getAllRoles(rolesPath string) ([]string, error) {
	// Get absolute path
	absPath, err := filepath.Abs(rolesPath)
	if err != nil {
		return nil, fmt.Errorf("error resolving absolute path: %v", err)
	}
	rp.Log("Scanning roles directory: %s\n", absPath)

	// Check if directory exists
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("error accessing roles directory: %v", err)
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("roles path is not a directory")
	}

	// List all entries in the directory
	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, fmt.Errorf("error reading roles irectory: %v", err)
	}

	rp.Log("Found %d entries in roles directory\n", len(entries))

	// Process entries to find roles
	var roles []string
	for _, entry := range entries {
		rp.processDirectoryEntry(entry, absPath, &roles)
	}

	// Print debug info if no roles found
	if len(roles) == 0 {
		rp.Log("WARNING: No roles found in directory %s\n", absPath)
		rp.printDirectoryDebugInfo(entries, absPath)
	}

	rp.Log("Found %d roles in directory\n", len(roles))
	return roles, nil
}
