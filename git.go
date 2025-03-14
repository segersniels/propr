package main

import (
	"os/exec"
	"sort"
	"strings"

	"github.com/charmbracelet/log"
)

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(stdout)), nil
}

func getAllBranches() ([]string, error) {
	// First, make sure we have the latest remote branches
	fetchCmd := exec.Command("git", "fetch", "--prune")
	if err := fetchCmd.Run(); err != nil {
		log.Debug("Failed to fetch remote branches", "error", err)
	}

	// Get local branches
	localCmd := exec.Command("git", "branch", "--format=%(refname:short)")
	localOut, err := localCmd.Output()
	if err != nil {
		return nil, err
	}

	// Get remote branches (excluding HEAD)
	remoteCmd := exec.Command("git", "branch", "-r", "--format=%(refname:short)")
	remoteOut, err := remoteCmd.Output()
	if err != nil {
		return nil, err
	}

	// Process local branches
	localBranches := strings.Split(strings.TrimSpace(string(localOut)), "\n")

	// Process remote branches and remove the "origin/" prefix
	remoteBranches := strings.Split(strings.TrimSpace(string(remoteOut)), "\n")
	var cleanedRemoteBranches []string
	for _, branch := range remoteBranches {
		if branch != "" && !strings.Contains(branch, "HEAD") {
			// Remove the "origin/" prefix for display
			parts := strings.SplitN(branch, "/", 2)
			if len(parts) > 1 {
				cleanedRemoteBranches = append(cleanedRemoteBranches, parts[1])
			}
		}
	}

	// Combine and deduplicate branches
	allBranches := make(map[string]bool)
	for _, branch := range localBranches {
		if branch != "" {
			allBranches[branch] = true
		}
	}

	for _, branch := range cleanedRemoteBranches {
		if branch != "" {
			allBranches[branch] = true
		}
	}

	// Convert map to slice
	var result []string
	for branch := range allBranches {
		result = append(result, branch)
	}

	// Sort branches alphabetically for better UX
	sort.Strings(result)

	return result, nil
}

func getDiff(current string, branch string) (string, error) {
	cmd := exec.Command("git", "diff", "origin/"+branch+".."+current)
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(stdout), nil
}

func fetchRemoteOrigin() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(stdout)), nil
}

type RepositoryInformation struct {
	Name  string
	Owner string
}

func getRepositoryInformation() (*RepositoryInformation, error) {
	origin, err := fetchRemoteOrigin()
	if err != nil {
		return nil, err
	}

	split := strings.Split(origin, "/")
	name := strings.TrimSuffix(split[len(split)-1], ".git")
	owner := split[len(split)-2]

	return &RepositoryInformation{
		Name:  name,
		Owner: owner,
	}, nil
}

func getCommitMessages(current string, branch string) ([]string, error) {
	cmd := exec.Command("git", "log", "--oneline", "origin/"+branch+".."+current)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(string(stdout), "\n"), nil
}
