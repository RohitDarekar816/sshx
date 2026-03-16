import os
import sys
import logging
import re
import hashlib
import json
from typing import Dict, Iterable, List, Optional, Set, Tuple

import requests

SONAR_TOKEN = os.getenv("SONAR_TOKEN")
GITHUB_TOKEN = os.getenv("GIT_ACCESS_TOKEN")
REPO = os.getenv("GIT_REPO", "owner/repo")

SONAR_URL = "https://sonarcloud.io/api/issues/search"
COMPONENT_KEY = "RohitDarekar816_sshx"
SONAR_ORG = os.getenv("SONAR_ORG")


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
    organization: Optional[str],
) -> Dict[str, object]:
    if rule_key in cache:
        return cache[rule_key]
    if not organization:
        logging.warning("Skipping rule fetch for %s because SONAR_ORG is not set.", rule_key)
        cache[rule_key] = {}
        return cache[rule_key]
    params = {"key": rule_key, "organization": organization}
    resp = session.get(
        "https://sonarcloud.io/api/rules/show",
        params=params,
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


def issue_fingerprint(component: str, line: Optional[object], rule_key: str, message: str) -> str:
    raw = f"{component}|{line}|{rule_key}|{message}"
    return hashlib.sha256(raw.encode("utf-8")).hexdigest()[:16]


def load_existing_issue_markers(
    session: requests.Session,
    repo: str,
) -> Tuple[Set[str], Set[str]]:
    sonar_keys: Set[str] = set()
    fingerprints: Set[str] = set()

    cache_dir = os.getenv("XDG_CACHE_HOME", "/tmp")
    cache_path = os.path.join(cache_dir, "sonar_to_github_issues_cache.json")
    cache_etag: Optional[str] = None
    cache_data: Optional[dict] = None

    try:
        with open(cache_path, "r", encoding="utf-8") as handle:
            cache_data = json.load(handle)
            cache_etag = cache_data.get("etag")
    except FileNotFoundError:
        cache_data = None
    except (json.JSONDecodeError, OSError):
        cache_data = None

    page = 1
    used_cache = False
    while True:
        headers = {}
        if page == 1 and cache_etag:
            headers["If-None-Match"] = cache_etag
        resp = session.get(
            f"https://api.github.com/repos/{repo}/issues",
            params={"state": "all", "per_page": 100, "page": page},
            headers=headers,
            timeout=20,
        )
        if resp.status_code == 304 and page == 1 and cache_data:
            cached_keys = cache_data.get("sonar_keys", [])
            cached_fps = cache_data.get("fingerprints", [])
            if isinstance(cached_keys, list):
                sonar_keys.update(k for k in cached_keys if isinstance(k, str))
            if isinstance(cached_fps, list):
                fingerprints.update(k for k in cached_fps if isinstance(k, str))
            used_cache = True
            break
        if not resp.ok:
            raise RuntimeError(
                f"GitHub issues list failed: {resp.status_code} {resp.text[:200]}"
            )
        issues = resp.json()
        if not issues:
            break
        for issue in issues:
            body = str(issue.get("body", "") or "")
            match_key = re.search(r"Sonar Issue Key:\s*([A-Za-z0-9:_-]+)", body)
            if match_key:
                sonar_keys.add(match_key.group(1))
            match_fp = re.search(r"Sonar Fingerprint:\s*([0-9a-f]{16})", body)
            if match_fp:
                fingerprints.add(match_fp.group(1))
        page += 1

    if not used_cache and cache_path:
        new_cache = {
            "etag": resp.headers.get("ETag") if "resp" in locals() else None,
            "sonar_keys": sorted(sonar_keys),
            "fingerprints": sorted(fingerprints),
        }
        try:
            os.makedirs(os.path.dirname(cache_path), exist_ok=True)
            with open(cache_path, "w", encoding="utf-8") as handle:
                json.dump(new_cache, handle)
        except OSError:
            pass

    return sonar_keys, fingerprints


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

    sonar_keys, fingerprints = load_existing_issue_markers(session, REPO)

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

            if sonar_issue_key in sonar_keys or fingerprint in fingerprints:
                logging.info("Skipping existing issue: %s", sonar_issue_key)
                continue

            label = issue_label(severity)
            labels = [label, "sonarqube"] if label != "sonarqube" else ["sonarqube"]

            rule_name = ""
            rule_desc = ""
            rule_why = ""
            rule_fix = ""
            if rule_key:
                rule = sonar_rule_details(sonar_session, rule_key, rule_cache, SONAR_ORG)
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
