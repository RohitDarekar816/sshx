# sshx

Git-backed SSH profile manager for teams.

sshx stores server profiles in a Git repository so teams can share and version-control SSH access. It provides a fast CLI to add, list, edit, remove, and connect to servers, with encrypted passwords, encrypted private keys, and TOTP-based MFA.

---

## Why sshx

- **Git-native**: all changes are committed and pushed.
- **Team-friendly**: share profiles through a repo instead of local configs.
- **Fast CLI**: connect with `sshx <server>` or `sshx connect <server>`.
- **Secure options**: encrypt passwords and private keys at rest with MFA for unlock.

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

This creates or clones the repo, registers you in `users.json`, and sets up TOTP MFA.

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
Authenticate with the Git repository that stores profiles. This also sets up TOTP MFA on first run.

### `sshx add [name]`
Add a new server profile.

Common flags:

- `--host` (required)
- `--user` (default: `root`)
- `--port` (default: `22`)
- `--key` (default: `~/.ssh/id_rsa`)
- `--key-ref` (reference to an encrypted key stored in the repo)
- `--password` (optional, stored encrypted)

### `sshx list`
List all servers in the repo.

### `sshx connect [server]`
Connect to a server by name.

Flags:

- No flags. Prompts for vault passphrase and TOTP code.

### `sshx [server]`
Direct shortcut to connect.

### `sshx edit [server]`
Edit server fields.

Flags:

- `--host`, `--user`, `--port`, `--key`, `--key-ref`
- `--password` to set an encrypted password
- `--clear-password` to remove stored password
- `--clear-key-ref` to remove stored key reference

### `sshx remove [server]`
Remove a server profile.

### `sshx key add|list|remove`
Manage encrypted private keys stored in the repo.

Example:

```bash
./sshx key add prod --file ~/.ssh/id_rsa
./sshx key list
```

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
key_ref: prod
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

## Encrypted Key Storage

Private keys can be stored in the Git repo encrypted at rest. Use `sshx key add` to encrypt and store a key, then reference it from a server profile with `--key-ref`.

Keys are stored under:

```
keys/
```

---

## MFA (TOTP)

TOTP is required for sensitive operations and for connecting to servers. On first `sshx auth`, the CLI shows a QR code in the terminal, plus the otpauth URL and manual secret for fallback.

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
- Sync and conflict resolution helper
- Profile schema validation and auto-repair

See `ROADMAP.md` for details.

---

## Validation and Auto-Repair

sshx validates server profiles on load and before saving. It auto-repairs common issues like missing defaults (user, port) and inconsistent auth fields.

---

## Conflict Resolution

When pushing changes, sshx detects non-fast-forward errors, performs a pull with rebase, and retries the push automatically.

---

## Contributing

Please read `CONTRIBUTING.md` before opening a PR.

---

## License

MIT
