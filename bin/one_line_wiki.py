#!/usr/bin/env - python3
import argparse
import json
import os
import sys
from typing import Dict, List

sys.stdout.reconfigure(line_buffering=True)

p = argparse.ArgumentParser()
args = p.parse_args()


def get_one_line_wikis(json_obj: Dict) -> List:
    lines = []
    profile = json_obj["profile"]
    for wiki in json_obj["wikis"]:
        line = []
        line.append(profile)
        line.append(str(wiki["id"]))
        line.append(str(wiki["name"]))
        lines.append(":".join(line))
    return lines


try:
    line = sys.stdin.readline()
    while line:
        line = line.strip("\n")
        if len(line) > 0:
            lines = get_one_line_wikis(json.loads(line))
            print("\n".join(lines))
            line = sys.stdin.readline()
except BrokenPipeError:
    devnull = os.open(os.devnull, os.O_WRONLY)
    os.dup2(devnull, sys.stdout.fileno())
    sys.exit(1)
