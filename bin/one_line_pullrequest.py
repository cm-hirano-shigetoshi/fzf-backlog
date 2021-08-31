#!/usr/bin/env - python3
import argparse
import json
import os
import sys
from typing import Dict, List

sys.stdout.reconfigure(line_buffering=True)

p = argparse.ArgumentParser()
args = p.parse_args()


def coloring(status: str) -> str:
    if status == "Open":
        return "\033[31m" + status + "\033[0m"
    elif status == "Merged":
        return "\033[32m" + status + "\033[0m"
    elif status == "Closed":
        return "\033[90m" + status + "\033[0m"
    return status


def get_one_line_pullrequests(json_obj: Dict) -> List:
    lines = []
    profile = json_obj["profile"]
    if not json_obj["repositories"]:
        return lines
    for repo in json_obj["repositories"]:
        if not repo["pullRequests"]:
            continue
        for pullrequest in repo["pullRequests"]:
            line = []
            line.append(profile)
            line.append(repo["repositoryName"])
            line.append("#" + str(pullrequest["number"]))
            line.append(coloring(pullrequest["status"]["name"]))
            if pullrequest["assignee"]:
                line.append(pullrequest["assignee"]["name"])
            line.append(pullrequest["summary"])
            lines.append(":".join(line))
    return lines


try:
    line = sys.stdin.readline()
    while line:
        line = line.strip("\n")
        if len(line) > 0:
            lines = get_one_line_pullrequests(json.loads(line))
            if len(lines) > 0:
                print("\n".join(lines))
            line = sys.stdin.readline()
except BrokenPipeError:
    devnull = os.open(os.devnull, os.O_WRONLY)
    os.dup2(devnull, sys.stdout.fileno())
    sys.exit(1)
