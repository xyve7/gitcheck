package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func traverse_directories(directories *[]string, path string, user string) {
	// Get the entries
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Print(err)
		return
	}

	// If the directory is empty
	if len(entries) == 2 {
		return
	}

	for _, entry := range entries {
		// If its the git directories, get the remote repositories
		if entry.IsDir() && entry.Name() == ".git" {
			data, err := exec.Command("git", "-C", path, "remote", "-v").Output()
			if err != nil {
				log.Print(err)
				continue
			}
			// We only care about the fetch
			remotes := strings.Split(string(data), "\n")
			if len(remotes) < 1 {
				continue
			}
			split_fetch := strings.Split(remotes[0], "/")
			if len(split_fetch) < 4 {
				continue
			}
			// Check if its the right user
			if user != split_fetch[3] {
				continue
			}

			// Check for uncommited changes
			commit, err := exec.Command("git", "-C", path, "status").Output()
			if err != nil {
				log.Print(err)
				continue
			}
			if strings.Contains(string(commit), "Changes not staged for commit") {
				*directories = append(*directories, path)
			}

		} else if entry.IsDir() {
			// Recursively search
			traverse_directories(directories, filepath.Join(path, entry.Name()), user)
		}
	}
}

func main() {
	// Set the values which will be checked
	var user string
	var root string
	flag.StringVar(&user, "user", "", "user which should be searched")
	flag.StringVar(&root, "root", "", "from where the .git directories should be searched")
	flag.Parse()

	if len(user) == 0 {
		log.Fatal("No user provided")
	}
	if len(root) == 0 {
		log.Fatal("No root path provided")
	}

	git_directories := make([]string, 0)
	traverse_directories(&git_directories, root, user)

	for _, dir := range git_directories {
		fmt.Printf("uncommited changes: %s\n", dir)
	}
}
