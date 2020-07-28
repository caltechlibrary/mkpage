#!/usr/bin/env python3

import os
import sys
import json

def filter_authors(src, limit, suffix):
    l = json.loads(src)
    # Find our authors list and tweak it.
    names = [];
    for name in l:
        s = f"{name['family']}, {name['given']}"
        names.append(s)
    if (limit > 0) and len(l) > limit:
        names = names[0:limit]
        if (suffix != ''):
            names.append(suffix)
    src = json.dumps(names)
    print(src)

if __name__ == "__main__":
    # get the file to process
    args = sys.argv[1:]
    if len(args) == 0:
        print("ERROR: nothing to process")
        sys.exit(1)
    f_name = ''
    limit = 0
    suffix = ''
    if len(args) > 0:
        f_name = args[0]
    if len(args) > 1:
        limit = int(args[1])
    if len(args) > 2:
        suffix = ' '.join(args[2:])
        if suffix.startswith('"') or suffix.startswith("'"):
            suffix = suffix[1:len(suffix)-1]
    with open(f_name) as f:
        src = f.read()
        filter_authors(src, limit, suffix)
