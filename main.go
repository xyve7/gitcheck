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
			// Check if the fetch contains the username
			contains := strings.Contains(remotes[0], user)
			if !contains {
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
	// Get user and root directory
	var user string
	var root string
	flag.StringVar(&user, "user", "", "user which should be searched")
	flag.StringVar(&root, "root", "", "from where the .git directories should be searched")
	flag.Parse()

	// Check if the values were provided
	if len(user) == 0 {
		log.Fatal("No user provided")
	}
	if len(root) == 0 {
		log.Fatal("No root path provided")
	}

	// Create the slice which will store the directories
	git_directories := make([]string, 0)
	traverse_directories(&git_directories, root, user)

	// Print out the uncommited changes
	fmt.Println("Uncomitted changes:")
	for _, dir := range git_directories {
		fmt.Printf("\t%s\n", dir)
	}
}
