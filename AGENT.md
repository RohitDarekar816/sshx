# AGENT.md

## Project Name

sshx

## Project Description

sshx is a CLI tool written in Go that allows users to manage SSH server profiles stored in a Git repository.

The tool is designed to provide a terminal-based experience similar to SSH profile managers, enabling users to quickly connect to servers without manually remembering hostnames, usernames, or SSH configurations.

Instead of storing profiles locally, sshx stores all configuration in a Git repository. This allows teams to share and version-control SSH access configurations.

The CLI will allow users to:

- Authenticate with a Git repository
- Register users
- Store SSH server profiles
- List available servers
- Connect to servers via SSH
- Synchronize configuration through Git

The long-term goal is to create a Git-backed SSH profile manager suitable for DevOps teams.

---

# Current Development Status

The project is currently in **Phase 2** of development.

Completed features:

- CLI framework using Cobra
- Git repository authentication
- User registration and storage
- Server profile creation
- Server listing
- Server connection

Next planned phase:

- Direct connection command (`sshx <server-name>`)
- Server removal
- YAML configuration
- Team access controls
- SSH key management

---

# CLI Commands Implemented

## sshx auth

Purpose:
Authenticate sshx with a Git repository that stores SSH profiles.

Example:

sshx auth [git@github.com](mailto:git@github.com):USER/sshx-profiles.git

Steps performed:

1. Initialize local sshx directory
2. Clone the Git repository
3. If repository is empty:
   - Initialize repository
   - Create users.json
   - Commit and push

4. Prompt user for:
   - name
   - email

5. Register user in users.json if not already present
6. Commit and push changes

---

## sshx add

Purpose:
Add a new SSH server profile.

Example:

sshx add prod --host 1.2.3.4 --user ubuntu

Flags:

--host server hostname or IP (required)
--user ssh username (default: root)
--port ssh port (default: 22)
--key ssh private key path (default: ~/.ssh/id_rsa)

Behavior:

1. Create server configuration
2. Save configuration in repository
3. Commit changes
4. Push to remote Git repository

Server profiles are currently stored as JSON.

Example stored file:

servers/prod.json

Example content:

{
"name": "prod",
"host": "1.2.3.4",
"user": "ubuntu",
"port": 22,
"key": "~/.ssh/id_rsa"
}

---

## sshx list

Purpose:
Display all registered servers.

Example:

sshx list

Output example:

Servers:

- prod (ubuntu@1.2.3.4)
- staging (ec2-user@2.2.2.2)

The command reads all server files from:

servers/

---

## sshx connect

Purpose:
Connect to a server using stored configuration.

Example:

sshx connect prod

Behavior:

1. Load server configuration
2. Construct SSH command
3. Execute system ssh client

Example internal command executed:

ssh ubuntu@1.2.3.4

---

# Repository Structure

sshx (CLI source code)

cmd/
root.go
auth.go
add.go
list.go
connect.go

internal/

config/
config.go

git/
git.go
repo.go

server/
server.go
loader.go

ssh/
(connect logic may be moved here later)

user/
user.go

main.go

---

# Remote Git Repository Structure

The Git repository used by sshx stores shared configuration.

Example structure:

sshx-profiles/

users.json

servers/
prod.json
staging.json
db.json

---

# Local Directory Structure

Local sshx data is stored under the user's home directory.

Example:

~/.sshx/

repo/
(cloned git repository)

---

# Key Design Concepts

## Git-Based Storage

All configuration is stored in a Git repository to provide:

- version control
- team collaboration
- centralized configuration
- audit history

Every configuration change performs:

1. file write
2. git commit
3. git push

---

## CLI Framework

The CLI is built using the Cobra framework.

Each command is implemented as a separate file under:

cmd/

---

## Internal Package Design

internal/config
Handles directory paths and local configuration.

internal/git
Handles Git operations:

- clone
- initialize
- commit
- push

internal/server
Handles server configuration:

- server struct
- save server
- load server
- list servers

internal/user
Handles users.json operations.

---

# Current Limitations

The project currently has the following limitations:

1. Server profiles stored in JSON instead of YAML
2. No direct command like `sshx prod`
3. No permission control
4. No SSH key distribution
5. No server removal command

---

# Planned Features (Next Phases)

Phase 3

Add command shortcut:

sshx prod

This will automatically resolve server name and connect.

---

Phase 4

Server management improvements:

sshx remove <server>
sshx edit <server>

---

Phase 5

Switch configuration format from JSON to YAML.

Example:

servers/prod.yaml

name: prod
host: 1.2.3.4
user: ubuntu
port: 22
key: ~/.ssh/id_rsa

---

Phase 6

Team access management:

groups/
devops.json

This will allow controlling which users can access specific servers.

---

Phase 7

SSH key management:

sshx push-key prod

This will automatically install public keys on servers.

---

# Development Guidelines

When modifying this project:

1. Maintain modular architecture inside internal/
2. Avoid placing business logic inside cmd/
3. Ensure all changes commit to Git repository
4. Prefer simple and readable code over complex abstractions
5. Keep CLI user experience fast and minimal

---

# Target Users

Primary users are:

- DevOps engineers
- Cloud engineers
- SRE teams
- developers managing multiple servers

---

# Long Term Vision

sshx aims to become a lightweight, Git-powered SSH access manager for teams.

The goal is to allow teams to manage infrastructure access using simple Git workflows and CLI commands.

---

# Maintainer

Project created and maintained by:

Rohit Darekar

GitHub repository:
https://github.com/RohitDarekar816/sshx
