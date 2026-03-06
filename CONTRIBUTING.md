# CONTRIBUTING.md

Thank you for your interest in contributing to sshx.

This project aims to become a lightweight Git-powered SSH access manager for teams.

---

# Development Setup

Requirements:

- Go 1.22 or newer
- Git installed
- SSH client available

Clone the repository:

git clone https://github.com/RohitDarekar816/sshx.git

Enter project directory:

cd sshx

Build the project:

go build -o sshx

---

# Running the CLI

Example usage:

sshx auth [git@github.com](mailto:git@github.com):USER/sshx-profiles.git

Add a server:

sshx add prod --host 1.2.3.4 --user ubuntu

List servers:

sshx list

Connect:

sshx connect prod

---

# Project Structure

cmd/

Contains CLI commands.

Each command has its own file.

internal/

Contains business logic packages.

Packages include:

config
git
server
user
ssh

main.go

Entry point of the application.

---

# Coding Guidelines

Follow these guidelines when contributing:

1. Keep command logic minimal in cmd/
2. Implement business logic in internal packages
3. Write clear and readable Go code
4. Avoid unnecessary abstractions
5. Maintain modular structure

---

# Commit Guidelines

Use clear commit messages.

Examples:

Add server command

Fix repo initialization bug

Improve user registration

---

# Pull Requests

When submitting a pull request:

1. Describe the change clearly
2. Reference related issues if any
3. Ensure the code builds successfully

---

# Feature Requests

Feature suggestions are welcome.

Please open a GitHub issue describing:

- The problem
- The proposed solution
- Example CLI usage

---

# Code of Conduct

Be respectful and constructive in discussions.

We welcome contributions from developers of all experience levels.

---

# Maintainer

Rohit Darekar

GitHub:
https://github.com/RohitDarekar816/sshx
