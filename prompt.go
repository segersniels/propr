package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
)

const SYSTEM_MESSAGE = `You are responsible to write a concise GitHub PR description.
You will be provided with the current branch name, the diff and the commit messages to provide you with enough
context to write a proper description. Analyze the code changes and provide a concise explanation of the changes,
their context and why they were made. If the provided message is not a diff respond with an appropriate message.
Only answer with the raw markdown description matching the template, do not include any other text.
Don't wrap your response in a markdown code block since GitHub will render it properly.
`

var FILES_TO_IGNORE = []string{
	"package-lock.json",
	"yarn.lock",
	"npm-debug.log",
	"yarn-debug.log",
	"yarn-error.log",
	".pnpm-debug.log",
	"Cargo.lock",
	"Gemfile.lock",
	"mix.lock",
	"Pipfile.lock",
	"composer.lock",
	"go.sum",
	"bun.lockb",
	"bun.lock",
}

func splitDiffIntoChunks(diff string) []string {
	split := strings.Split(diff, "diff --git")[1:]
	for i, chunk := range split {
		split[i] = strings.TrimSpace(chunk)
	}

	return split
}

func removeLockFiles(chunks []string) []string {
	var wg sync.WaitGroup

	filtered := make(chan string)

	for _, chunk := range chunks {
		wg.Add(1)

		go func(chunk string) {
			defer wg.Done()
			shouldIgnore := false
			header := strings.Split(chunk, "\n")[0]

			// Check if the first line contains any of the files to ignore
			for _, file := range FILES_TO_IGNORE {
				if strings.Contains(header, file) {
					log.Debug("Ignoring", "file", file)
					shouldIgnore = true
				}
			}

			if !shouldIgnore {
				log.Debug("Adding", "header", header)
				filtered <- chunk
			}
		}(chunk)
	}

	go func() {
		wg.Wait()
		close(filtered)
	}()

	var result []string
	for chunk := range filtered {
		result = append(result, chunk)
	}

	return result
}

// Split the diff in chunks and remove any lock files to save on tokens
func prepareDiff(diff string) string {
	chunks := splitDiffIntoChunks(diff)

	return strings.Join(removeLockFiles(chunks), "\n")
}

func generateSystemMessageForDiff(systemMessage string, template string) string {
	return fmt.Sprintf("%s\n\nFollow this exact template to write your description:\n\n%s", systemMessage, template)
}
