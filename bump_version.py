import subprocess

result = subprocess.run(["git", "log", "-1", "--pretty=%s"], capture_output=True, text=True)
commit = result.stdout.strip()

with open("VERSION", "r") as fp:
    version = fp.read()
    print(version)

[major, minor, patch] = map(int, version.split("."))

if "BREAKING CHANGE" in commit or commit.startswith("feat!:"):
    major += 1
    minor = 0
    patch = 0
elif commit.startswith("feat:"):
    minor += 1
    patch = 0
elif commit.startswith("fix:"):
    patch += 1

with open("VERSION", "w") as fp:
    fp.write(f"{major}.{minor}.{patch}")