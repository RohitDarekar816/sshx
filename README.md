# sshx

Git-backed SSH profile manager for teams.

sshx stores server profiles in a Git repository so teams can share and version-control SSH access. It provides a fast CLI to add, list, edit, remove, and connect to servers, with optional encrypted passwords.

---

## Why sshx

- **Git-native**: all changes are committed and pushed.
- **Team-friendly**: share profiles through a repo instead of local configs.
- **Fast CLI**: connect with `sshx <server>` or `sshx connect <server>`.
- **Secure options**: encrypt passwords at rest; supports SSH keys by default.

---

## Quick Start

### 1) Install

```bash
go build -o sshx
```

### 2) Authenticate with a repo

```bash
./sshx auth git@github.com:USER/sshx-profiles.git
```

This creates or clones the repo and registers you in `users.json`.

### 3) Add a server

```bash
./sshx add prod --host 1.2.3.4 --user ubuntu
```

### 4) List servers

```bash
./sshx list
```

### 5) Connect

```bash
./sshx prod
```

---

## Commands

### `sshx auth [repo-url]`
Authenticate with the Git repository that stores profiles.

### `sshx add [name]`
Add a new server profile.

Common flags:

- `--host` (required)
- `--user` (default: `root`)
- `--port` (default: `22`)
- `--key` (default: `~/.ssh/id_rsa`)
- `--password` (optional, stored encrypted)

### `sshx list`
List all servers in the repo.

### `sshx connect [server]`
Connect to a server by name.

### `sshx [server]`
Direct shortcut to connect.

### `sshx edit [server]`
Edit server fields.

Flags:

- `--host`, `--user`, `--port`, `--key`
- `--password` to set an encrypted password
- `--clear-password` to remove stored password

### `sshx remove [server]`
Remove a server profile.

---

## Server Profile Format

Server profiles live in the Git repo under `servers/` and are stored as YAML.

Example:

```yaml
name: prod
host: 1.2.3.4
user: ubuntu
port: 22
key: ~/.ssh/id_rsa
```

sshx can still read existing JSON profiles for backward compatibility.

---

## Password Storage

If you use `--password`, sshx encrypts the password locally using AES-256-GCM with a passphrase. The encrypted value is stored in the server profile.

To connect using a stored password, `sshpass` must be installed:

```bash
sudo apt install sshpass
```

---

## Local Data

sshx stores local data under:

```
~/.sshx/
```

The cloned repo lives at:

```
~/.sshx/repo/
```

---

## Requirements

- Go 1.24+
- Git
- SSH client
- `sshpass` (optional, only if you use encrypted passwords)

---

## Roadmap

Next planned features:

- Group-based access control
- SSH key distribution
- Cloud provider imports

See `ROADMAP.md` for details.

---

## Contributing

Please read `CONTRIBUTING.md` before opening a PR.

---

## License

MIT
