---
allowed-tools:
  - Bash:*
  - Read
argument-hint: "[version-type]"
description: "Auto-merge iteration to main with version bump and tag - user pushes manually"
---

# Iteration Commit Workflow

You are an iteration commit orchestrator. Your task is to squash merge the current iteration branch into main, create a synthesized commit message, and add a semantic version tag. The user will push manually.

## Process

### Step 1: Get Current State

Execute these commands sequentially to understand the current state. Use simple commands to avoid bash parsing issues:

```bash
# Get current branch name
git branch --show-current
```

Store the branch name, then check if it's main (if so, error and exit). If branch matches pattern "iteration-N", extract the iteration number and show iteration details:

```bash
tm iteration show <ITERATION_NUM> --full
```

Get latest version tag:

```bash
git describe --tags --abbrev=0 2>/dev/null
```

If no tags exist, assume v0.0.0.

Show what will be merged (replace `<BRANCH>` with actual branch name):

```bash
echo "--- Commits to be squashed ---"
git log main..<BRANCH> --oneline --no-decorate | head -20
```

```bash
echo "--- Files changed ---"
git diff main..<BRANCH> --stat | head -30
```

```bash
echo "--- Key changes summary ---"
git diff main..<BRANCH> --name-only | head -50
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

If user provided version type in argument (major/minor/patch), use that. Otherwise, **automatically detect** based on changes.

**Analysis Commands** (replace `<BRANCH>` with actual branch name):

Get commit messages:
```bash
git log main..<BRANCH> --pretty=format:"%s"
```

Get changed files:
```bash
git diff main..<BRANCH> --name-status
```

Count commits:
```bash
git log main..<BRANCH> --oneline | wc -l
```

Count files:
```bash
git diff main..<BRANCH> --name-only | wc -l
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

Execute immediately (no confirmation needed). Replace `<BRANCH>`, `<VERSION>`, `<COMMIT_MESSAGE>`, and `<TAG_MESSAGE>` with actual values:

Switch to main and update:
```bash
git checkout main && git pull origin main
```

Squash merge the iteration branch:
```bash
git merge --squash <BRANCH>
```

Create commit with synthesized message (use HEREDOC for proper formatting):
```bash
git commit -m "$(cat <<'EOF'
<COMMIT_MESSAGE>
EOF
)"
```

Create annotated tag (use HEREDOC for proper formatting):
```bash
git tag -a <VERSION> -m "$(cat <<'EOF'
<TAG_MESSAGE>
EOF
)"
```

Show result:
```bash
git log -1 --oneline && echo && git tag -l <VERSION> -n99
```

### Step 5: Final Summary

Display final summary (replace placeholders with actual values):

```
=== ITERATION COMMIT COMPLETE ===

✅ Squash merged: <BRANCH> → main
✅ Commit created: <first line of commit message>
✅ Tag created: <VERSION>

Current branch: main
Original branch: <BRANCH> (preserved - delete manually if needed)

Next steps:
1. Review the commit and tag
2. Push to origin:
   git push origin main
   git push origin <VERSION>
3. Verify changes on GitHub/remote
4. Optionally delete branch: git branch -d <BRANCH>
```

## Success Criteria

- ✅ Current branch squash merged into main
- ✅ Synthesized commit message (not individual commits)
- ✅ Semantic version tag created with proper message
- ✅ User has clear summary of what happened
- ✅ Push instructions provided to user
- ✅ Original branch preserved for manual cleanup

## Error Handling

- If already on main: Error and exit
- If merge conflicts: Abort merge, show conflicts, ask user to resolve manually
- If commit fails: Show error and abort

## Important Notes

- **No automatic push**: User must manually push to origin
- **Branch preservation**: Original branch is NOT deleted (user decides)
- **Bash simplicity**: Use simple commands to avoid eval parsing issues
- **HEREDOC for messages**: Always use HEREDOC syntax for multi-line commit/tag messages
- **Placeholders**: Replace `<BRANCH>`, `<VERSION>`, etc. with actual values in commands

Remember: This is a fully automated workflow command. Execute all steps automatically without confirmations, and provide clear feedback throughout.
