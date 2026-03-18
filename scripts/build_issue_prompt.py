#!/usr/bin/env python3
import argparse
import json
import os
import sys
import urllib.parse
import urllib.request


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Build a prompt from a GitHub issue for the Codex CLI."
    )
    parser.add_argument(
        "--repo",
        required=True,
        help="GitHub repo in owner/name format for fetching an issue.",
    )
    parser.add_argument(
        "--issue-number",
        type=int,
        help="Specific GitHub issue number to use.",
    )
    parser.add_argument(
        "--issue-pick",
        choices=["oldest", "newest"],
        default="oldest",
        help="How to pick an issue when --issue-number is not provided.",
    )
    parser.add_argument(
        "--label",
        help="Optional label to filter issues by.",
    )
    parser.add_argument(
        "--state",
        choices=["open", "closed", "all"],
        default="open",
        help="Issue state to query.",
    )
    return parser.parse_args()


def build_issue_prompt(repo: str, issue: dict) -> str:
    number = issue.get("number", "unknown")
    title = issue.get("title") or "(no title)"
    body = (issue.get("body") or "").strip()
    if not body:
        body = "(no description provided)"

    return (
        "You are an AI coding agent. Fix the following GitHub issue and create a PR.\n\n"
        f"Repository: {repo}\n"
        f"Issue #{number}: {title}\n\n"
        "Issue description:\n"
        f"{body}\n\n"
        "Requirements:\n"
        "- Implement a correct fix in the codebase.\n"
        "- Add or update tests as needed.\n"
        "- Use the GitHub CLI (`gh`) to open a PR targeting the develop branch (not master), with a clear summary and reference the issue number.\n"
    )


def github_api_request(url: str, token: str | None) -> dict | list:
    headers = {
        "Accept": "application/vnd.github+json",
        "User-Agent": "codex-gha-runner",
    }
    if token:
        headers["Authorization"] = f"Bearer {token}"

    req = urllib.request.Request(url, headers=headers)
    with urllib.request.urlopen(req, timeout=30) as resp:
        data = resp.read()
    return json.loads(data.decode("utf-8"))


def fetch_issue(repo: str, issue_number: int | None, label: str | None, state: str, issue_pick: str) -> dict:
    token = os.environ.get("GITHUB_TOKEN")
    if issue_number is not None:
        url = f"https://api.github.com/repos/{repo}/issues/{issue_number}"
        return github_api_request(url, token)

    direction = "asc" if issue_pick == "oldest" else "desc"
    params = {
        "state": state,
        "per_page": "30",
        "sort": "created",
        "direction": direction,
    }
    if label:
        params["labels"] = label

    query = urllib.parse.urlencode(params)
    url = f"https://api.github.com/repos/{repo}/issues?{query}"
    issues = github_api_request(url, token)
    if not isinstance(issues, list) or not issues:
        raise RuntimeError("No issues found with the given filters.")

    if state == "open":
        issues = [i for i in issues if "pull_request" not in i]

    if not issues:
        raise RuntimeError("No non-PR issues found with the given filters.")

    return issues[0]


def main() -> int:
    args = parse_args()

    issue = fetch_issue(
        repo=args.repo,
        issue_number=args.issue_number,
        label=args.label,
        state=args.state,
        issue_pick=args.issue_pick,
    )

    prompt = build_issue_prompt(args.repo, issue)
    sys.stdout.write(prompt)
    return 0


if __name__ == "__main__":
    sys.exit(main())
