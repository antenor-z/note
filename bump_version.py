import subprocess
import re

result = subprocess.run(["git", "log", "-1", "--pretty=%B"], capture_output=True, text=True)
commit = result.stdout.strip()

with open("VERSION", "r") as fp:
    version = fp.read().strip()
    print(version)

[major, minor, patch] = map(int, version.split("."))

if "BREAKING CHANGE" in commit or \
    re.search(r"feat(\([a-zA-Z0-9]+\))?!:", commit) or \
    re.search(r"fix(\([a-zA-Z0-9]+\))?!:", commit):
    major += 1
    minor = 0
    patch = 0
elif re.search(r"feat(\([a-zA-Z0-9]+\))?:", commit):
    minor += 1
    patch = 0
elif re.search(r"fix(\([a-zA-Z0-9]+\))?:", commit):
    patch += 1

new_version = f"{major}.{minor}.{patch}"

if new_version != version:
    print(f"Bump version: {version} to {new_version}")
    with open("VERSION", "w") as fp:
        fp.write(f"{major}.{minor}.{patch}")