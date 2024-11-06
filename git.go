package main

import (
	"os/exec"
	"strings"
)

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(stdout)), nil
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
