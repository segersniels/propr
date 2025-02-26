# Propr

[![npm](https://img.shields.io/npm/v/@segersniels/propr)](https://www.npmjs.com/package/@segersniels/propr)![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/segersniels/propr/ci.yml)

Generate GitHub PR descriptions from the command line with the help of AI.
`propr` aims to populate a basic PR description right from your terminal so you can focus on more important things.

## Install

### NPM

```bash
npm install -g @segersniels/propr
```

### Script

```bash
# Install in the current directory
curl -sSL https://raw.githubusercontent.com/segersniels/propr/master/scripts/install.sh | bash
# Install in /usr/local/bin
curl -sSL https://raw.githubusercontent.com/segersniels/propr/master/scripts/install.sh | sudo bash -s /usr/local/bin
```

### Manual

1. Download the latest binary from the [releases](https://github.com/segersniels/propr/releases/latest) page for your system
2. Rename the binary to `propr`
3. Copy the binary to a location in your `$PATH`

## Prerequisites

Before using Propr, you'll need to set up the following:

1. **GitHub Token**: Set the `GITHUB_TOKEN` environment variable with a valid GitHub token that has permissions to create pull requests.

   ```bash
   export GITHUB_TOKEN=your_github_token
   ```

2. **AI Provider API Key**: Depending on which AI model you want to use, set one of the following environment variables:

   - For OpenAI models (default): `export OPENAI_API_KEY=your_openai_api_key`
   - For Anthropic models: `export ANTHROPIC_API_KEY=your_anthropic_api_key`
   - For DeepSeek models: `export DEEPSEEK_API_KEY=your_deepseek_api_key`

## Usage

Before you can get started, write some code and push it to a branch. Then depending on your needs, you can use the following commands:

```bash
NAME:
   propr - Generate your PRs from the command line with AI

USAGE:
   propr [global options] command [command options]

COMMANDS:
   create    Creates a PR with a generated description
   generate  Generates a PR description and outputs it
   config    Configure propr to your liking
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Configuration Options

Propr can be configured with the following options:

1. **Model**: Choose which AI model to use for generating PR descriptions.
2. **Prompt**: Customize the system prompt used to generate PR descriptions.
3. **Template**: Define the template structure for your PR descriptions (default: `# Description`).
4. **Pretty Print**: Enable or disable pretty printing of the generated output (default: `true`).

Example configuration:

```json
{
  "model": "gpt-4o-mini",
  "prompt": "You are responsible to write a concise GitHub PR description...",
  "template": "# Description\n\n## Changes\n\n## Why",
  "pretty_print": true
}
```

## Advanced Usage

### Customizing Templates

You can customize the PR description template to match your team's standards. For example:

```bash
propr config init
```

Then, when prompted, enter a template like:

```
# Description

## Changes
-

## Impact
-

## Testing
-
```

### Debug Mode

Enable debug mode to see more detailed logs:

```bash
DEBUG=true propr generate
```

## Common Gotchas and Troubleshooting

### API Key Issues

If you encounter errors like `OPENAI_API_KEY is not set`, make sure you've set the appropriate environment variable for your chosen model.

### Git Repository Issues

1. **No Remote Origin**: Propr requires a GitHub repository with a remote origin set up. Make sure your local repository has a remote origin pointing to GitHub.
2. **No Changes**: If you get an error like `not enough changes found to generate`, make sure you have committed and pushed changes to your branch.
3. **Branch Comparison**: By default, Propr compares your current branch to the repository's default branch. If you want to compare against a different branch, use the `--branch` flag.
4. **Repository Format**: Propr expects the remote origin URL to be in a standard GitHub format (e.g., `https://github.com/username/repo.git` or `git@github.com:username/repo.git`).

### Rate Limiting

If you encounter rate limiting issues with the GitHub API, consider using a personal access token with higher rate limits.

### Large Diffs

For very large changes, the AI model might struggle to generate a comprehensive description due to token limits. In such cases:

- Consider breaking your PR into smaller, more focused changes
- Propr automatically filters out lock files like `package-lock.json`, `yarn.lock`, etc. to save on tokens

## Environment Variables

- `GITHUB_TOKEN`: Required for creating PRs and accessing repository information
- `OPENAI_API_KEY`: Required when using OpenAI models
- `ANTHROPIC_API_KEY`: Required when using Anthropic models
- `DEEPSEEK_API_KEY`: Required when using DeepSeek models
- `DEBUG`: Set to any value to enable debug logging

## How It Works

Propr works by:

1. Fetching the diff between your current branch and the target branch
2. Collecting commit messages from your branch
3. Sending this information to the configured AI model
4. Generating a PR description based on the changes
5. Optionally creating a PR with the generated description

The tool automatically filters out lock files and other large generated files to optimize token usage.

## Contributing

Contributions are welcome! Check out the [GitHub repository](https://github.com/segersniels/propr) for more information.
