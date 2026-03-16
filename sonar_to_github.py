import os
import sys
import logging
import re
import hashlib
from typing import Dict, Iterable, List, Optional, Set, Tuple

import requests

SONAR_TOKEN = os.getenv("SONAR_TOKEN")
GITHUB_TOKEN = os.getenv("GIT_ACCESS_TOKEN")
REPO = os.getenv("GIT_REPO", "owner/repo")

SONAR_URL = "https://sonarcloud.io/api/issues/search"
COMPONENT_KEY = "RohitDarekar816_sshx"


def require_env(name: str) -> str:
    value = os.getenv(name)
    if not value:
        print(f"Missing required environment variable: {name}", file=sys.stderr)
        sys.exit(1)
    return value


def sonarqube_issues(
    sonar_token: str,
    component_key: str,
    page_size: int = 100,
) -> Iterable[Dict[str, object]]:
    page = 1
    total = None
    session = requests.Session()
    session.auth = (sonar_token, "")

    while total is None or (page - 1) * page_size < total:
        params = {
            "componentKeys": component_key,
            "resolved": "false",
            "p": page,
            "ps": page_size,
        }
        resp = session.get(SONAR_URL, params=params, timeout=20)
        if not resp.ok:
            raise RuntimeError(
                f"Sonar API failed: {resp.status_code} {resp.text[:200]}"
            )
        data = resp.json()
        issues = data.get("issues", [])
        for issue in issues:
            yield issue
        paging = data.get("paging", {})
        total = paging.get("total", total if total is not None else 0)
        page += 1


def sonar_rule_details(
    session: requests.Session,
    rule_key: str,
    cache: Dict[str, Dict[str, object]],
) -> Dict[str, object]:
    if rule_key in cache:
        return cache[rule_key]
    resp = session.get(
        "https://sonarcloud.io/api/rules/show",
        params={"key": rule_key},
        timeout=20,
    )
    if not resp.ok:
        raise RuntimeError(
            f"Sonar rule fetch failed: {resp.status_code} {resp.text[:200]}"
        )
    rule = resp.json().get("rule", {})
    cache[rule_key] = rule
    return rule


def issue_label(severity: str) -> str:
    label_map = {
        "BLOCKER": "critical",
        "CRITICAL": "high",
        "MAJOR": "medium",
        "MINOR": "low",
        "INFO": "info",
    }
    return label_map.get(severity, "sonarqube")


def normalize_component(component: str) -> str:
    # Sonar component format: "key:path/to/file"
    if ":" in component:
        return component.split(":", 1)[1]
    return component


def strip_html(text: str) -> str:
    return re.sub(r"<[^>]+>", "", text)


def truncate(text: str, max_len: int = 1200) -> str:
    if len(text) <= max_len:
        return text
    return text[: max_len - 3].rstrip() + "..."


def extract_rule_sections(rule: Dict[str, object]) -> Tuple[Optional[str], Optional[str]]:
    sections = rule.get("descriptionSections", [])
    why = None
    fix = None
    if isinstance(sections, list):
        for section in sections:
            key = str(section.get("key", "")).lower()
            content = str(section.get("content", "")).strip()
            if not content:
                continue
            if key in {"root_cause", "impact", "why"} and not why:
                why = content
            if key in {"how_to_fix", "how_to_fix_java", "how_to_fix_cs", "fix"} and not fix:
                fix = content
    return why, fix


def github_session(github_token: str) -> requests.Session:
    session = requests.Session()
    session.headers.update(
        {
            "Authorization": f"token {github_token}",
            "Accept": "application/vnd.github+json",
        }
    )
    return session


def ensure_labels(session: requests.Session, repo: str, labels: List[str]) -> None:
    existing: Set[str] = set()
    page = 1
    while True:
        resp = session.get(
            f"https://api.github.com/repos/{repo}/labels",
            params={"per_page": 100, "page": page},
            timeout=20,
        )
        if not resp.ok:
            raise RuntimeError(
                f"GitHub labels list failed: {resp.status_code} {resp.text[:200]}"
            )
        data = resp.json()
        if not data:
            break
        for item in data:
            name = item.get("name")
            if isinstance(name, str):
                existing.add(name)
        page += 1

    label_colors = {
        "critical": "b60205",
        "high": "d93f0b",
        "medium": "fbca04",
        "low": "0e8a16",
        "info": "1d76db",
        "sonarqube": "5319e7",
    }

    for label in labels:
        if label in existing:
            continue
        color = label_colors.get(label, "ededed")
        resp = session.post(
            f"https://api.github.com/repos/{repo}/labels",
            json={"name": label, "color": color},
            timeout=20,
        )
        if not resp.ok:
            raise RuntimeError(
                f"GitHub label create failed: {resp.status_code} {resp.text[:200]}"
            )


def search_issues(session: requests.Session, query: str) -> bool:
    resp = session.get(
        "https://api.github.com/search/issues",
        params={"q": query, "per_page": 1},
        timeout=20,
    )
    if not resp.ok:
        raise RuntimeError(
            f"GitHub search failed: {resp.status_code} {resp.text[:200]}"
        )
    data = resp.json()
    return bool(data.get("total_count", 0))


def sanitize_query_term(value: str, max_len: int = 80) -> str:
    value = value.replace('"', " ").strip()
    if len(value) > max_len:
        value = value[:max_len].rstrip()
    return value


def issue_fingerprint(component: str, line: Optional[object], rule_key: str, message: str) -> str:
    raw = f"{component}|{line}|{rule_key}|{message}"
    return hashlib.sha256(raw.encode("utf-8")).hexdigest()[:16]


def github_issue_exists(
    session: requests.Session,
    repo: str,
    sonar_issue_key: str,
    fingerprint: str,
    file_path: str,
    line_display: object,
    message: str,
) -> bool:
    if sonar_issue_key:
        query = f'repo:{repo} type:issue in:body "{sonar_issue_key}"'
        if search_issues(session, query):
            return True

    query = f'repo:{repo} type:issue in:body "Sonar Fingerprint: {fingerprint}"'
    if search_issues(session, query):
        return True

    msg = sanitize_query_term(message)
    file_term = sanitize_query_term(file_path)
    line_term = sanitize_query_term(str(line_display))
    if msg and file_term and line_term:
        query = (
            f'repo:{repo} type:issue in:body "File: {file_term}" '
            f'"Line: {line_term}" "{msg}"'
        )
        if search_issues(session, query):
            return True

    if msg:
        query = f'repo:{repo} type:issue in:title "{msg}"'
        if search_issues(session, query):
            return True

    return False


def create_github_issue(
    session: requests.Session,
    repo: str,
    title: str,
    body: str,
    labels: List[str],
) -> None:
    resp = session.post(
        f"https://api.github.com/repos/{repo}/issues",
        json={"title": title, "body": body, "labels": labels},
        timeout=20,
    )
    if not resp.ok:
        raise RuntimeError(
            f"GitHub issue create failed: {resp.status_code} {resp.text[:200]}"
        )


def main() -> None:
    logging.basicConfig(
        level=logging.INFO,
        format="%(levelname)s: %(message)s",
    )
    sonar_token = require_env("SONAR_TOKEN")
    github_token = require_env("GIT_ACCESS_TOKEN")

    session = github_session(github_token)

    labels_to_ensure = {"sonarqube", "critical", "high", "medium", "low", "info"}
    ensure_labels(session, REPO, sorted(labels_to_ensure))

    sonar_session = requests.Session()
    sonar_session.auth = (sonar_token, "")
    rule_cache: Dict[str, Dict[str, object]] = {}

    for issue in sonarqube_issues(sonar_token, COMPONENT_KEY):
        try:
            severity = str(issue.get("severity", "INFO"))
            message = str(issue.get("message", "")).strip()
            component = str(issue.get("component", ""))
            line = issue.get("line")
            line_display = line if line is not None else "N/A"
            sonar_issue_key = str(issue.get("key", "")).strip()
            rule_key = str(issue.get("rule", "")).strip()
            issue_type = str(issue.get("type", "CODE_SMELL"))
            project_key = str(issue.get("project", COMPONENT_KEY))
            fingerprint = issue_fingerprint(component, line, rule_key, message)
            file_path = normalize_component(component)

            if not sonar_issue_key:
                logging.warning("Skipping issue without key: %s", message)
                continue

            if github_issue_exists(
                session,
                REPO,
                sonar_issue_key,
                fingerprint,
                file_path,
                line_display,
                message,
            ):
                logging.info("Skipping existing issue: %s", sonar_issue_key)
                continue

            label = issue_label(severity)
            labels = [label, "sonarqube"] if label != "sonarqube" else ["sonarqube"]

            rule_name = ""
            rule_desc = ""
            rule_why = ""
            rule_fix = ""
            if rule_key:
                rule = sonar_rule_details(sonar_session, rule_key, rule_cache)
                rule_name = str(rule.get("name", "")).strip()
                rule_desc = str(rule.get("mdDesc") or rule.get("htmlDesc") or "").strip()
                rule_desc = strip_html(rule_desc)
                rule_why, rule_fix = extract_rule_sections(rule)

            body = (
                "SonarQube Issue Detected\n\n"
                f"File: {file_path}\n"
                f"Line: {line_display}\n"
                f"Severity: {severity}\n\n"
                f"Type: {issue_type}\n"
                f"Rule Key: {rule_key}\n"
                f"Rule Name: {rule_name}\n\n"
                f"Sonar Issue Key: {sonar_issue_key}\n\n"
                f"Sonar Fingerprint: {fingerprint}\n\n"
                f"Issue URL: https://sonarcloud.io/project/issues?open={sonar_issue_key}&id={project_key}\n\n"
                "Message:\n"
                f"{message}\n\n"
                "Why this is an issue:\n"
                f"{truncate(rule_why or rule_desc or 'Not provided by Sonar.')}\n\n"
                "How to fix:\n"
                f"{truncate(rule_fix or 'Refer to the rule guidance in SonarCloud.')}\n"
            )

            # GitHub title limit is 256 chars.
            title = f"[SonarQube] {message}"
            if len(title) > 256:
                title = title[:253] + "..."

            create_github_issue(session, REPO, title, body, labels)
            logging.info("Created issue: %s", sonar_issue_key)
        except Exception as exc:
            logging.error("Failed to process issue: %s", exc)


if __name__ == "__main__":
    main()
