# ROADMAP.md

This document outlines the planned development roadmap for sshx.

---

# Phase 1 (Completed)

CLI foundation.

Features implemented:

- Cobra CLI framework
- Project structure
- Git repository authentication
- User registration
- users.json management

Command implemented:

sshx auth

---

# Phase 2 (Completed)

Basic server management.

Features implemented:

- Add server profile
- List server profiles
- Connect to server

Commands implemented:

sshx add
sshx list
sshx connect

Server configuration stored in:

servers/

---

# Phase 3 (Completed)

Improve CLI experience.

Planned feature:

Direct server command.

Example:

sshx prod

Instead of:

sshx connect prod

Implementation idea:

Unknown command fallback will attempt to resolve server name.

---

# Phase 4 (Completed)

Server management improvements.

New commands:

sshx remove <server>

sshx edit <server>

This allows modifying existing server configurations.

---

# Phase 5 (In Progress)

Configuration format improvements.

Change server profiles from JSON to YAML.

Status:

- YAML read/write support in CLI
- Backward-compatible loading of existing JSON profiles

Example:

servers/prod.yaml

YAML is easier for humans to read and edit.

---

# Phase 6

Access control system.

Add group-based permissions.

Repository structure:

groups/

Example:

groups/devops.json

Example permissions:

devops team can access prod servers.

---

# Phase 7

SSH key management.

New command:

sshx push-key prod

This will automatically install a user's public key on the server.

---

# Phase 8

Cloud integration.

Ability to import servers from cloud providers.

Potential commands:

sshx import aws
sshx import gcp
sshx import azure

This would automatically register servers from cloud infrastructure.

---

# Long-Term Vision

sshx aims to become a lightweight Git-powered access management tool for infrastructure teams.

Key goals:

- Git-based access control
- Simple CLI experience
- DevOps-friendly workflow
- Team collaboration
