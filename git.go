package main

import (
	"bufio"
	"os/exec"
	"strings"
)

func getDefaultBranch() (string, error) {
	cmd := exec.Command("git", "remote", "show", "origin")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse the output to get the default branch
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "HEAD branch") {
			parts := strings.Split(line, ":")
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return string(stdout), nil
}

func getDiff(branch string) (string, error) {
	cmd := exec.Command("git", "diff", "origin/"+branch)
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
