// package utils ...
package utils

import (
	"fmt"
	"go/build"
	"os"
)

func SetWorkingDir(path string) error {
	dir, err := importPathToDir(path)
	if err != nil {
		return fmt.Errorf("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		return fmt.Errorf("os.Chdir:", err)
	}
	return nil
}

// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}
