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
    if status == "未対応":
        return "\033[31m" + status + "\033[0m"
    elif status == "処理中":
        return "\033[34m" + status + "\033[0m"
    elif status == "処理済み":
        return "\033[32m" + status + "\033[0m"
    elif status == "完了":
        return "\033[33m" + status + "\033[0m"
    return status


def get_one_line_issues(json_obj: Dict) -> List:
    lines = []
    profile = json_obj["profile"]
    for issue in json_obj["issues"]:
        line = []
        line.append(profile)
        line.append(issue["issueKey"])
        line.append(coloring(issue["status"]["name"]))
        if issue["assignee"] is not None:
            line.append(issue["assignee"]["name"])
        else:
            line.append("")
        line.append(issue["summary"])
        lines.append(":".join(line))
    return lines


try:
    line = sys.stdin.readline()
    while line:
        line = line.strip("\n")
        if len(line) > 0:
            lines = get_one_line_issues(json.loads(line))
            print("\n".join(lines))
            line = sys.stdin.readline()
except BrokenPipeError:
    devnull = os.open(os.devnull, os.O_WRONLY)
    os.dup2(devnull, sys.stdout.fileno())
    sys.exit(1)
