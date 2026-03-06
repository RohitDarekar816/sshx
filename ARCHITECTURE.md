# ARCHITECTURE.md

## Overview

sshx is a Git-backed CLI tool for managing SSH server profiles.

Instead of storing SSH configurations locally, sshx stores them in a Git repository. This allows teams to share server profiles and manage infrastructure access through version control.

The CLI reads and writes configuration files in the repository and automatically commits and pushes changes.

---

# System Architecture

The architecture follows a modular Go design using internal packages.

The CLI layer handles user interaction, while internal packages handle logic such as Git operations, configuration management, and server handling.

High-level flow:

User CLI Command
↓
Cobra Command Handler (cmd/)
↓
Internal Packages (logic layer)
↓
Filesystem + Git repository
↓
SSH execution (system ssh)

---

# Component Overview

## CLI Layer

Directory:

cmd/

Responsible for:

- Parsing CLI commands
- Handling flags
- Calling internal services

Commands implemented:

auth
add
list
connect

Each command is implemented in its own file.

Example:

cmd/add.go
cmd/list.go

This keeps the CLI layer clean and modular.

---

# Internal Packages

## internal/config

Purpose:
Manage local configuration and filesystem paths.

Responsibilities:

- Initialize sshx base directory
- Determine repository directory
- Manage configuration paths

Example local structure:

~/.sshx/

repo/

The repo directory contains the cloned Git configuration repository.

---

## internal/git

Purpose:
Handle all Git operations.

Responsibilities:

- Clone repository
- Initialize empty repositories
- Commit changes
- Push changes

All configuration changes go through this module.

Typical workflow:

write file
git add
git commit
git push

The project uses the Go library:

go-git

---

## internal/user

Purpose:
Manage users.json file.

Responsibilities:

- Load users
- Check if user exists
- Add new user
- Save users file

users.json format example:

[
{
"name": "Rohit",
"email": "[rohit@example.com](mailto:rohit@example.com)"
}
]

This file allows tracking which users are registered in the sshx system.

---

## internal/server

Purpose:
Manage server profiles.

Responsibilities:

- Server struct definition
- Save server profiles
- Load server profiles
- List servers

Servers are stored in:

servers/

Example:

servers/prod.json

Example structure:

{
"name": "prod",
"host": "1.2.3.4",
"user": "ubuntu",
"port": 22,
"key": "~/.ssh/id_rsa"
}

---

## internal/ssh

Purpose:
Handle SSH connections.

Responsibilities:

- Construct SSH command
- Execute system ssh client
- Forward input/output streams

Example command executed:

ssh ubuntu@1.2.3.4

The tool intentionally relies on the system SSH client instead of reimplementing SSH.

---

# Data Flow Example

Example command:

sshx add prod --host 1.2.3.4 --user ubuntu

Flow:

CLI parses command
↓
cmd/add.go handles input
↓
internal/server.SaveServer() writes configuration file
↓
internal/git.CommitAndPush() commits change
↓
Git repository updated

---

# Repository Layout

sshx CLI source code:

cmd/
internal/
main.go

Git configuration repository:

users.json
servers/

Example:

sshx-profiles/

users.json

servers/
prod.json
staging.json

---

# Future Architecture Improvements

Planned improvements include:

1. YAML configuration instead of JSON
2. Permission control via groups
3. Direct command execution (`sshx prod`)
4. SSH key distribution
5. Cloud provider integration
6. Plugin system

---

# Design Principles

The project follows these design principles:

Simple CLI UX
Git-native configuration
Modular Go architecture
Minimal dependencies
Human-readable configuration
