#!/usr/bin/env python3
import argparse
import json
import os
import sys
import urllib.parse
import urllib.request

from oz_agent_sdk import OzAPI


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Create a new Oz agent run with a custom prompt."
    )
    parser.add_argument(
        "--prompt",
        help="Prompt to send to the Oz agent. If omitted, an issue will be fetched from GitHub.",
    )
    parser.add_argument(
        "--repo",
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
    parser.add_argument(
        "--environment-id",
        help="Optional Oz environment ID to run the agent in.",
    )
    parser.add_argument(
        "--model-id",
        help="Optional model ID to use for the run.",
    )
    parser.add_argument(
        "--base-prompt",
        help="Optional base prompt to prepend for the agent.",
    )
    parser.add_argument(
        "--name",
        help="Optional config name for traceability.",
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
        "- Open a PR with a clear summary and reference the issue number.\n"
    )


def github_api_request(url: str, token: str | None) -> dict | list:
    headers = {
        "Accept": "application/vnd.github+json",
        "User-Agent": "oz-agent-runner",
    }
    if token:
        headers["Authorization"] = f"Bearer {token}"

    req = urllib.request.Request(url, headers=headers)
    with urllib.request.urlopen(req, timeout=30) as resp:
        data = resp.read()
    return json.loads(data.decode("utf-8"))


def fetch_issue(repo: str, issue_number: int | None, label: str | None, state: str) -> dict:
    token = os.environ.get("GITHUB_TOKEN")
    if issue_number is not None:
        url = f"https://api.github.com/repos/{repo}/issues/{issue_number}"
        return github_api_request(url, token)

    params = {
        "state": state,
        "per_page": "30",
        "sort": "created",
        "direction": "asc",
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

    client = OzAPI(
        api_key=os.environ.get("WARP_API_KEY"),
    )

    prompt = args.prompt
    if not prompt:
        if not args.repo:
            raise SystemExit("--repo is required when --prompt is not provided.")

        issue = fetch_issue(
            repo=args.repo,
            issue_number=args.issue_number,
            label=args.label,
            state=args.state,
        )
        if args.issue_pick == "newest" and args.issue_number is None:
            # Re-fetch with descending sort to avoid extra list logic.
            token = os.environ.get("GITHUB_TOKEN")
            params = {
                "state": args.state,
                "per_page": "30",
                "sort": "created",
                "direction": "desc",
            }
            if args.label:
                params["labels"] = args.label
            query = urllib.parse.urlencode(params)
            url = f"https://api.github.com/repos/{args.repo}/issues?{query}"
            issues = github_api_request(url, token)
            if isinstance(issues, list) and issues:
                issue = issues[0]

        prompt = build_issue_prompt(args.repo, issue)

    config = {}
    if args.environment_id:
        config["environment_id"] = args.environment_id
    if args.model_id:
        config["model_id"] = args.model_id
    if args.base_prompt:
        config["base_prompt"] = args.base_prompt
    if args.name:
        config["name"] = args.name

    response = client.agent.run(
        prompt=prompt,
        config=config or None,
    )

    print(response.run_id)
    return 0


if __name__ == "__main__":
    sys.exit(main())
