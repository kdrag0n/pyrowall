// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
var Default = Build

func getGitCommit() (string, error) {
	return sh.Output("git", "rev-list", "-1", "HEAD")
}

// Build builds an executable.
func Build() error {
	mg.Deps(Deps)
	fmt.Println("Building...")

	gitCommit, err := getGitCommit()
	if err != nil {
		return err
	}

	return sh.Run("go", "build", "-ldflags", "-X main.GitCommit="+gitCommit, ".")
}

// BuildStripped builds an executables with debugging information and symbol tables stripped.
func BuildStripped() error {
	mg.Deps(Deps)
	fmt.Println("Building stripped version...")

	gitCommit, err := getGitCommit()
	if err != nil {
		return err
	}

	return sh.Run("go", "build", "-ldflags", "-s -w -X main.GitCommit="+gitCommit, ".")
}

// Deps downloads the project's dependencies.
func Deps() error {
	fmt.Println("Downloading dependencies...")
	return sh.Run("go", "mod", "download")
}

// Clean removes the built executable.
func Clean() {
	fmt.Println("Cleaning...")
	os.Remove("pyrowall")
}
