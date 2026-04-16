---
description: Build, release, and publish to Homebrew
argument-hint: <patch|minor|major>
---

# Release

Perform a full release: version bump, build, GitHub release, and Homebrew tap update.

## Arguments

- `$ARGUMENTS` — version bump type: `patch`, `minor`, or `major`

## Instructions

You are performing a release of the `lin` CLI. Follow these steps exactly.

### Pre-flight

1. Confirm the working tree is clean (`git st`). If not, stop and ask.
2. Run `make test` and `go vet ./...`. If either fails, stop and fix.
3. Determine the current version from the latest git tag (`git describe --tags --abbrev=0`) and show what bump will happen. If no tag exists, start at `0.1.0`.

### Step 1: Version bump, tag, and push

Calculate the new version by bumping the current tag:

```bash
# Get current version
current=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "0.0.0")
# Split into parts and bump based on type
IFS='.' read -r major minor patch <<< "$current"
```

Apply the bump type ($ARGUMENTS):
- `patch`: increment patch
- `minor`: increment minor, reset patch to 0
- `major`: increment major, reset minor and patch to 0

Then tag and push:

```bash
git tag "v${new_version}"
git push origin main "v${new_version}"
```

### Step 2: Cross-compile

```bash
rm -rf dist/
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.version=${new_version}" -o "dist/lin-darwin-arm64" ./cmd/lin
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=${new_version}" -o "dist/lin-darwin-amd64" ./cmd/lin
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=${new_version}" -o "dist/lin-linux-amd64" ./cmd/lin
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X main.version=${new_version}" -o "dist/lin-linux-arm64" ./cmd/lin
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=${new_version}" -o "dist/lin-windows-amd64.exe" ./cmd/lin

# Create tarballs + checksums
cd dist
for bin in lin-darwin-arm64 lin-darwin-amd64 lin-linux-amd64 lin-linux-arm64; do
  tar czf "${bin}.tar.gz" "$bin"
done
shasum -a 256 *.tar.gz lin-windows-amd64.exe > checksums-sha256.txt
cd ..
```

Test the native binary before proceeding:

```bash
./dist/lin-darwin-arm64 --version
./dist/lin-darwin-arm64 user me
```

### Step 3: Create GitHub release

```bash
prev_tag=$(git tag --sort=-v:refname | head -2 | tail -1)
notes=$(git log --pretty=format:"- %s" "${prev_tag}..v${new_version}" --no-merges | grep -v "^- v[0-9]")

gh release create "v${new_version}" dist/*.tar.gz dist/lin-windows-amd64.exe dist/checksums-sha256.txt \
  --title "v${new_version}" \
  --notes "$notes"
```

Verify: `gh release view "v${new_version}"`

### Step 4: Update Homebrew tap

The Homebrew formula lives in `../homebrew-tap` relative to this repo's root.

```bash
ls ../homebrew-tap/Formula/lin.rb
```

Read checksums from `dist/checksums-sha256.txt` and update the formula:

1. Read `../homebrew-tap/Formula/lin.rb`
2. Update version, URLs (use `v${new_version}`), SHA256 values, assert_match version
3. Commit and push:
   ```bash
   cd ../homebrew-tap
   git add Formula/lin.rb
   git commit -m "lin ${new_version}"
   git push
   cd -
   ```

**IMPORTANT:** Always `cd` back to the lin repo after updating the tap.

### Step 5: Report

Show the user:

- New version number
- GitHub release URL
- Homebrew tap commit
- `brew install shhac/tap/lin` command for new users
- `brew upgrade shhac/tap/lin` command for existing users
