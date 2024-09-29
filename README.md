# Propr

![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/segersniels/propr/ci.yml)

Generate GitHub PR descriptions from the command line with the help of AI.
`propr` aims to populate a basic PR description right from your terminal so you can focus on more important things.

<p align="center">
<img src="./resources/logo.png" width="250">

## Install

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

## Usage

```
NAME:
   propr - Generate your PRs from the command line with AI

USAGE:
   propr [global options] command [command options]

VERSION:
   x.x.x

COMMANDS:
   create    Creates a PR with a generated description
   generate  Generates a PR description and outputs it
   config    Configure propr to your liking
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```