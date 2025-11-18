---
allowed-tools:
  - Bash:*
  - Read
argument-hint: "[version-type] [--no-push]"
description: "Auto-merge iteration to main with version bump - no confirmations, branch preserved"
---

# Iteration Commit Workflow

You are an iteration commit orchestrator. Your task is to squash merge the current iteration branch into main, create a synthesized commit message, add a semantic version tag, and push to origin.

## Process

### Step 1: Get Current State

Execute these commands to understand the current state:

```bash
# Get current branch name
CURRENT_BRANCH=$(git branch --show-current)
echo "Current branch: $CURRENT_BRANCH"

# Check if on main (error if true)
if [ "$CURRENT_BRANCH" = "main" ]; then
    echo "ERROR: Already on main branch. Switch to iteration branch first."
    exit 1
fi

# Get current iteration info (if branch starts with "iteration-")
if [[ "$CURRENT_BRANCH" =~ ^iteration-([0-9]+) ]]; then
    ITERATION_NUM="${BASH_REMATCH[1]}"
    echo "Iteration number: $ITERATION_NUM"
    dw task-manager iteration show $ITERATION_NUM --full
fi

# Get latest version tag
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo "Latest version: $LATEST_TAG"

# Show what will be merged
echo -e "\n--- Commits to be squashed ---"
git log main..$CURRENT_BRANCH --oneline --no-decorate | head -20

echo -e "\n--- Files changed ---"
git diff main..$CURRENT_BRANCH --stat | head -30

echo -e "\n--- Key changes summary ---"
git diff main..$CURRENT_BRANCH --name-only | head -50
```

### Step 2: Analyze Changes and Create Synthesized Commit Message

Based on the output above, create a **synthesized, high-level commit message** that:
- Follows conventional commits format: `type(scope): summary`
- Focuses on WHAT changed (user-facing features/fixes)
- Groups related changes into categories
- Keeps it concise (3-5 bullets max per section)
- References the iteration number

**DO NOT** copy all individual commits. **DO** synthesize the overall purpose and impact.

**Template:**
```
type(scope): brief summary of iteration's main accomplishment

Major changes:
- Feature 1
- Feature 2
- Feature 3

UX/Bug fixes:
- Improvement 1
- Fix 1

Closes #<iteration-number> - <iteration name>
```

**Common types:**
- `feat:` - New features
- `fix:` - Bug fixes
- `refactor:` - Code restructuring
- `docs:` - Documentation
- `style:` - UI/formatting changes
- `test:` - Test additions

### Step 3: Automatically Detect Version Bump

**Semantic Versioning: MAJOR.MINOR.PATCH**

- **MAJOR** (1.x.x → 2.x.x): Breaking changes, incompatible API changes
- **MINOR** (x.1.x → x.2.x): New features (backward-compatible)
- **PATCH** (x.x.1 → x.x.2): Bug fixes only (backward-compatible)

If user provided version type in $1 (major/minor/patch), use that. Otherwise, **automatically detect** based on changes:

**Analysis Commands:**
```bash
# Get commit messages to analyze
git log main..HEAD --pretty=format:"%s" > /tmp/commits.txt

# Get changed files
git diff main..HEAD --name-status > /tmp/changes.txt

# Count changes
TOTAL_COMMITS=$(git log main..HEAD --oneline | wc -l)
TOTAL_FILES=$(git diff main..HEAD --name-only | wc -l)
```

**MAJOR bump indicators** (breaking changes - highest priority):
- Commit messages contain: `BREAKING CHANGE:`, `breaking:`, or `!:` after type (e.g., `feat!:`)
- Migration files added (e.g., `migrations/`, `schema/`)
- Major version in package files updated
- Config file format changes
- Removed commands or APIs
- Changed database schema incompatibly

**MINOR bump indicators** (new features):
- Commit types: `feat:`, `feature:`
- New files in significant locations:
  - New `cmd/` directories
  - New `.claude/commands/` files
  - New `pkg/` packages
  - New UI components/views
- Commit messages contain: "add", "new", "implement", "introduce"
- Iteration numbers (indicates feature work)
- More than 5 files changed

**PATCH bump indicators** (bug fixes only):
- Commit types: `fix:`, `bugfix:`, `hotfix:`
- Only `*_test.go` files changed
- Only documentation files changed (`*.md`, `docs/`)
- Commit messages contain: "fix", "patch", "correct", "repair"
- Fewer than 5 files changed
- No new files, only modifications

**Detection Algorithm:**

```
1. Check for MAJOR indicators → If found: MAJOR
2. Check for MINOR indicators → If found: MINOR
3. Default: PATCH
```

**Display Detection Result:**
```
=== VERSION BUMP DETECTION ===

Analyzed: [TOTAL_COMMITS] commits, [TOTAL_FILES] files changed

Detected bump: [MAJOR/MINOR/PATCH]

Key indicators:
- [Indicator 1 that triggered this decision]
- [Indicator 2]
- [Indicator 3]

Current version: [LATEST_TAG]
New version: [NEW_VERSION]
```

Calculate new version based on detection.

### Step 4: Execute Merge and Tag

Execute immediately (no confirmation needed):

```bash
# Store current branch for later
CURRENT_BRANCH=$(git branch --show-current)

# Switch to main and update
git checkout main
git pull origin main

# Squash merge the iteration branch
git merge --squash $CURRENT_BRANCH

# Create commit with synthesized message
git commit -m "[COMMIT_MESSAGE]"

# Create annotated tag
git tag -a [NEW_VERSION] -m "[TAG_MESSAGE]"

# Show result
git log -1 --oneline
git tag -l [NEW_VERSION] -n99
```

### Step 5: Push to Origin

If --no-push flag is NOT provided:

```bash
# Push main and tags
git push origin main
git push origin [NEW_VERSION]

echo "✅ Pushed to origin: main and tag [NEW_VERSION]"
```

If --no-push flag IS provided:
```bash
echo "⚠️  Skipped push (--no-push flag). To push manually:"
echo "   git push origin main"
echo "   git push origin [NEW_VERSION]"
```

### Step 6: Final Summary

Display final summary:
```
=== ITERATION COMMIT COMPLETE ===

✅ Squash merged: $CURRENT_BRANCH → main
✅ Commit created: [first line of commit message]
✅ Tag created: [NEW_VERSION]
✅ Pushed to origin: [YES/NO]

Current branch: $CURRENT_BRANCH (kept - delete manually if needed)

Next steps:
- Verify changes on GitHub/remote
- Update CHANGELOG.md if you maintain one
- Create release notes from tag message
- Optionally delete branch: git branch -d $CURRENT_BRANCH
```

## Success Criteria

- ✅ Current branch squash merged into main
- ✅ Synthesized commit message (not individual commits)
- ✅ Semantic version tag created with proper message
- ✅ Changes pushed to origin (unless --no-push)
- ✅ User has clear summary of what happened
- ✅ Branch preserved for manual cleanup

## Error Handling

- If already on main: Error and exit
- If merge conflicts: Abort merge, show conflicts, ask user to resolve manually
- If git push fails: Show error, suggest manual push

Remember: This is a fully automated workflow command. Execute all steps automatically without confirmations, and provide clear feedback throughout.
