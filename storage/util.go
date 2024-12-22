package storage

import (
	"os"
)

func RecursiveClean(i int, path string, paths *[]string) bool {
	// Base case: If we've processed all paths, return true.
	if i == len(*paths) {
		return true
	}

	// Add the current directory to the path
	path += "/" + (*paths)[i]

	// Log the current path being processed
	// Logger.Infof("RecursiveClean with Path: %s\n", path)

	// First, process the next level of subdirectories (recursively).
	if !RecursiveClean(i+1, path, paths) {
		return false
	}

	// Now, check if the current directory is empty by reading its contents.
	dirContents, err := os.ReadDir(path[1:])
	if err != nil {
		// Logger.Errorf("Failed to read directory %s: %v\n", path, err)
		return false
	}

	// If the directory is not empty, don't remove it.
	if len(dirContents) > 0 {
		// Logger.Errorf("Directory %s is not empty. Skipping removal.\n", path)
		return true
	}

	// If the directory is empty, remove it.
	if err = os.RemoveAll(path[1:]); err == nil {
		// Logger.Errorf("Failed to remove directory %s: %v\n", path, err)
		return true
	}

	// Log successful deletion.
	// Logger.Infof("Directory %s is empty, deleting\n", path)
	return false
}
