---
allowed-tools:
  - Read
  - Edit
  - Task
  - TodoWrite
  - Bash(dw:*)
  - Bash(go:*)
  - Bash(go-arch-lint:*)
  - Bash(make:*)
  - Bash(git:*)
argument-hint: "[iteration-number-or-id]"
description: Implement TODO tasks in iteration with multi-agent workflow (plan → implement → verify). Ignores in-progress/review tasks.
---

# Iteration Implementation Command

You are an iteration orchestrator. Your task is to implement TODO tasks for an iteration using a multi-agent workflow: planning, implementation, and verification.

**Critical**:
- You ONLY implement tasks with status "todo"
- Tasks with status "in-progress", "review", "done", or "blocked" are IGNORED
- You orchestrate agents. You do NOT check acceptance criteria or close the iteration
- The user must verify all acceptance criteria and close the iteration manually

---

## Your Mission

Implement TODO tasks for an iteration from start to finish:
1. Identify target iteration (current, next by rank, or specified)
2. Filter to TODO tasks only (ignore "in-progress", "review", "done", "blocked")
3. Gather iteration context for TODO tasks (tasks, acceptance criteria, scope)
4. Launch planning agent (creates detailed implementation plan for TODO tasks)
5. Launch implementation agents (can run phases in parallel)
6. Launch verification agent (checks tests, linter, implementation quality)
7. Report completion and remind user to verify acceptance criteria

**CRITICAL**: Only implement tasks with status "todo". Leave all other tasks unchanged.

**You coordinate; sub-agents execute.**

---

## Process

### Phase 1: Identify Target Iteration

Execute these bash commands to determine the target iteration:

```bash
# Get current iteration
!`dw task-manager iteration current`
```

**Determine target**:

**If argument provided** (`$1`):
- If numeric (e.g., "3"): Use iteration number `$1`
- If starts with "TM-iter-": Use iteration ID `$1`
- Store as `TARGET_ITERATION`

**If no argument provided**:
- Parse `dw task-manager iteration current` output
- If current iteration exists and not completed: Use current iteration
- If no current or current completed:
  ```bash
  # List all iterations to find next by rank
  dw task-manager iteration list
  ```
  - Find first iteration with status "planned" or "in-progress" (sorted by rank)
  - Store as `TARGET_ITERATION`

**If no suitable iteration found**:
- Report: "No iterations ready to implement. Use `dw task-manager iteration create` to create one."
- Exit gracefully

---

### Phase 2: Start Iteration (if needed)

Check iteration status from Phase 1 output. If status is "planned", start it:

```bash
dw task-manager iteration start $TARGET_ITERATION
```

---

### Phase 2.5: Identify TODO Tasks

**CRITICAL**: This command only implements tasks with status "todo". Tasks with status "in-progress" or "review" are ignored.

Execute these commands to identify todo tasks:

```bash
# Get iteration details with task statuses
dw task-manager iteration show $TARGET_ITERATION --full

# Parse output to identify tasks with status "todo"
# Store list of TODO_TASKS (IDs of tasks with status "todo")
```

**Filter logic**:
- Include: Tasks with status "todo"
- Exclude: Tasks with status "in-progress", "review", "done", "blocked"

**If no todo tasks found**:
- Report: "No todo tasks in iteration [number]. All tasks are either in-progress, review, done, or blocked."
- List task statuses for user visibility
- Exit gracefully (no implementation needed)

**Continue only if**:
- At least one task has status "todo"
- Store TODO_TASKS list for planning agent

---

### Phase 3: Launch Planning Agent

**Create detailed planning prompt**:

```
You are the planning agent for iteration [iteration-number or iteration-id].

## Step 1: Explore Codebase Architecture

**Before planning**, use the Task tool with Explore agent to understand the codebase context.

Based on iteration scope (read iteration details first), explore relevant areas:

**For new features/commands**:
- Explore similar existing implementations
- Understand current patterns and conventions
- Find related domain models and repositories

**For refactoring/changes**:
- Explore the code being modified
- Understand current architecture and dependencies
- Find test patterns in use

**Example exploration**:
```
Use Task tool with:
- subagent_type: "Explore"
- description: "Explore [relevant area] architecture"
- prompt: "Explore the [package/area] to understand [what you need to know].

          Focus on:
          - [Specific aspect 1]
          - [Specific aspect 2]
          - [Specific aspect 3]

          Thoroughness: very thorough"
```

**Thoroughness guidance**:
- Simple iterations (1-2 small tasks): "quick" or "medium"
- Moderate iterations (3-4 tasks, new features): "medium" or "very thorough"
- Complex iterations (5+ tasks, architectural changes): "very thorough"

You may run **multiple explorations** for different areas if needed.

**Synthesize exploration findings** before creating the plan.

## Step 2: Gather Iteration Context

Execute these commands to get complete iteration information:

```bash
# Get complete iteration details with all tasks and full descriptions
dw task-manager iteration show [iteration-number] --full

# Get all acceptance criteria for all tasks in the iteration
dw task-manager ac list-iteration [iteration-number]
```

Parse the output to understand:
- Iteration name, goal, deliverable
- All tasks in scope (IDs, titles, descriptions, status)
- All acceptance criteria (descriptions, testing instructions, verification status)

**CRITICAL - TODO Tasks Only**:
- You are planning ONLY for tasks with status "todo"
- The following TODO task IDs must be implemented: [TODO_TASKS from Phase 2.5]
- Ignore tasks with status "in-progress", "review", "done", or "blocked"
- Focus your plan exclusively on the todo tasks listed above

## Step 3: Synthesize and Plan

Using **both exploration findings and iteration details**, create a comprehensive implementation plan that:

1. **Analyzes TODO tasks only** and their acceptance criteria
2. **Decomposes into implementation phases** (each phase should fit in agent context)
3. **Identifies dependencies** between phases
4. **Marks parallel opportunities** (which phases can run concurrently)
5. **Applies clean architecture principles** (follow DarwinFlow CLAUDE.md patterns, incorporate patterns from exploration)
6. **Ensures testability** (each phase should have clear verification steps)

## Architecture Context

- **Package structure**: See @docs/arch-index.md
- **Dependency rules**: pkg/pluginsdk imports nothing, internal/domain may import SDK only, etc.
- **Framework principle**: Framework is plugin-agnostic; plugin-specific types belong in plugin packages
- **Testing**: 70-80% coverage target, black-box testing (package pkgname_test)
- **Linter**: Zero violations required (go-arch-lint .)

## Plan Format

Return a detailed plan with:

### Phase 1: [Phase Name]
**Objective**: [What this phase accomplishes]
**Tasks involved**: [Which task IDs]
**Requirements**:
- [ ] Specific requirement 1
- [ ] Specific requirement 2
**Files to modify**: [Expected files]
**Testing strategy**: [How to verify]
**Can run in parallel with**: [None / Phase X, Phase Y]

### Phase 2: [Phase Name]
...

### Verification Phase: Final Checks
**Objective**: Verify all implementation is correct
**Requirements**:
- [ ] All tests pass (go test ./...)
- [ ] Zero linter violations (go-arch-lint .)
- [ ] All task acceptance criteria can be verified (manual check by user)
- [ ] Implementation matches plan
- [ ] No architectural violations

## Final Report Format

Return:
- Number of phases identified
- Which phases can run in parallel
- Critical dependencies to watch
- Estimated complexity (simple/moderate/complex)
- Any risks or concerns

Think deeply about:
- Clean architecture boundaries (informed by exploration)
- Existing patterns to follow (from exploration findings)
- Test strategy for each phase
- Minimal essential tests (don't over-test)
- Parallel execution opportunities
```

**Execute**:
```
Use Task tool with:
- subagent_type: "general-purpose" (requires design and planning)
- description: "Plan TODO tasks implementation"
- prompt: [constructed prompt above with TODO_TASKS list from Phase 2.5]
```

**IMPORTANT**: Replace `[TODO_TASKS from Phase 2.5]` in the prompt with the actual TODO task IDs identified in Phase 2.5.

**Review planning agent report**:
- Parse phases from report
- Note parallel execution opportunities
- Store phase details for next step
- Create TodoWrite with all phases

---

### Phase 4: Execute Implementation Phases

**For each phase in the plan** (respecting dependencies and parallel opportunities):

#### 4a. Select Sub-Agent Type

**Use junior-dev-executor when**:
- Phase has clear, well-specified requirements
- Implementation path is straightforward
- Can be executed directly without research

**Use general-purpose when**:
- Phase requires research or exploration
- Needs design decisions or trade-offs
- Requirements are less specific
- Involves discovery work

**Default**: If phase checklist is specific and actionable → junior-dev-executor. If requires discovery → general-purpose.

#### 4b. Construct Phase Prompt

```
You are implementing Phase [N]: [Phase Name] for iteration [Iteration Name]

## Iteration Context
[Brief 2-3 sentence summary of iteration goal]

## TODO Tasks Scope
**CRITICAL**: You are implementing ONLY tasks with status "todo".
- TODO task IDs in this iteration: [TODO_TASKS from Phase 2.5]
- Ignore tasks with status "in-progress", "review", "done", or "blocked"
- This phase addresses the following TODO tasks: [Task IDs for this phase]

## Phase Objective
[What this specific phase should accomplish]

## Phase Requirements
[Specific checklist items for this phase from planning agent]

## Related Tasks
[Task IDs and titles this phase addresses - must be from TODO_TASKS list]

## Acceptance Criteria Context
[Relevant AC from TODO tasks - for awareness, NOT for you to verify]
Note: You do NOT verify acceptance criteria. User will verify manually.

## Architecture Constraints
- Follow DarwinFlow package structure (see @docs/arch-index.md)
- Maintain dependency rules (SDK imports nothing, domain imports SDK only)
- Framework is plugin-agnostic
- Zero linter violations required
- Target 70-80% test coverage

## Verification
After implementation:
- [ ] Run tests: go test ./...
- [ ] Run linter: go-arch-lint .
- [ ] Update task status if applicable
- [ ] Note any blockers or issues

## Expected Deliverables
[What code, tests, or docs should be created/modified]

## Final Report Format
Return:
- What was implemented
- Files created/modified
- Test results (pass/fail with details)
- Linter results (violations count, details)
- Task status updates made
- Any issues or blockers
- Recommendations for next phase
```

#### 4c. Execute Phase

**Sequential execution** (default):
```
Use Task tool with:
- subagent_type: [junior-dev-executor OR general-purpose]
- description: "Implement [phase name]"
- prompt: [constructed prompt above]
- Wait for completion before next phase
```

**Parallel execution** (if planning agent identified):
```
If phases X, Y, Z can run in parallel:
- Launch multiple Task tools in single message
- Each with own phase prompt
- Wait for ALL to complete before dependent phases
```

#### 4d. Review Phase Report

**For each completed phase**:
- Mark TodoWrite item as completed
- Note files modified
- Check if tests passed
- Check if linter clean
- Update task statuses if mentioned in report

**If phase reports issues**:
- Do NOT mark todo as completed
- Do NOT proceed to dependent phases
- Report blocker to user
- Ask for guidance

**Update task statuses as appropriate**:
```bash
# If phase completes a TODO task, update status from "todo" to "done"
dw task-manager task update <task-id> --status done

# Note: Only update tasks that were in TODO_TASKS list
# Do NOT update tasks that were already "in-progress", "review", or "done"
```

---

### Phase 5: Launch Verification Agent

After all implementation phases complete, launch verification agent:

**Verification prompt**:

```
You are the verification agent for iteration [iteration-number].

## Implementation Context

**Implementation plan summary**: [Brief summary of phases from planning agent]
**Phases completed**: [List of phase names executed]
**Files modified**: [List key files from all phase reports]
**TODO tasks implemented**: [List of TODO_TASKS from Phase 2.5 that were implemented]

**CRITICAL - Scope of Verification**:
- This implementation focused ONLY on tasks with status "todo"
- Tasks with status "in-progress", "review", "done", or "blocked" were NOT modified
- Verify only the TODO tasks listed above were completed

## Step 1: Gather Current State

Execute these commands to verify current iteration state:

```bash
# Get current iteration status and task breakdown
dw task-manager iteration show [iteration-number] --full

# Get all acceptance criteria status
dw task-manager ac list-iteration [iteration-number]
```

**Filter verification scope**:
- Focus on tasks that were in TODO_TASKS list
- Confirm these tasks are now marked as "done" (if fully implemented)
- Confirm tasks that were "in-progress", "review", or "done" remain unchanged

## Step 2: Verification Checklist

### 1. Tests
- [ ] Run: go test ./...
- [ ] All tests pass
- [ ] No flaky tests
- [ ] Coverage is reasonable (70-80% target)

### 2. Architecture Linter
- [ ] Run: go-arch-lint .
- [ ] Zero violations
- [ ] No dependency rule violations
- [ ] Framework remains plugin-agnostic

### 3. Code Quality
- [ ] Read modified files (from implementation context above)
- [ ] Check clean architecture boundaries respected
- [ ] Verify proper error handling
- [ ] Confirm no code duplication
- [ ] Ensure proper separation of concerns

### 4. Implementation vs Plan
- [ ] Compare implementation to original plan
- [ ] All phase objectives met
- [ ] No missing functionality
- [ ] No scope creep or unrelated changes

### 5. Task Status Check
- [ ] Review task statuses from iteration show output
- [ ] All TODO tasks (from TODO_TASKS list) are now marked as "done"
- [ ] Tasks that were "in-progress", "review", or "done" remain unchanged
- [ ] No tasks stuck in wrong status

### 6. Acceptance Criteria Readiness
- [ ] Review AC for TODO tasks only (from ac list-iteration output)
- [ ] Confirm implementation enables user to verify each AC for TODO tasks
- [ ] Note: You do NOT verify AC yourself (user must do this)
- [ ] Check if testing instructions are clear for user
- [ ] AC for non-TODO tasks (in-progress/review/done) are not your concern

## Final Report Format

Return:

### Test Results
[Pass/fail, any failures, coverage notes]

### Linter Results
[Pass/fail, any violations]

### Code Quality Assessment
[Issues found, architectural concerns, clean code notes]

### Implementation Completeness
[Missing functionality, plan adherence, scope notes]

### Task Status Verification
[TODO tasks completed (from TODO_TASKS list), tasks unchanged (in-progress/review/done), status accuracy]

### Acceptance Criteria Readiness
[Are all AC for TODO tasks verifiable by user? Clear instructions?]

### Overall Assessment
- Ready for user acceptance: YES/NO
- Issues requiring immediate attention: [list]
- Recommendations for improvement: [list]
```

**Execute**:
```
Use Task tool with:
- subagent_type: "general-purpose" (requires analysis and verification)
- description: "Verify iteration implementation"
- prompt: [constructed prompt above]
```

**Review verification report**:
- Check if tests passed
- Check if linter clean
- Note any issues found
- Determine if ready for user acceptance

**If verification finds issues**:
- Create new phase: "Fix verification issues"
- Delegate to appropriate agent with issue details
- Run verification again
- Only proceed when clean

---

### Phase 6: Documentation Check

**Determine if documentation updates needed**:

- Did functionality change? → Check if README.md needs update
- Did workflow change? → Check if CLAUDE.md needs update
- Did architecture change? → Consider running `go-arch-lint docs`

**If updates needed**:
- Note for user (don't auto-update without user confirmation for major docs)
- Or make updates if changes are minor/obvious

---

### Phase 7: Create Git Commit

**MANDATORY**: After all implementation and verification is complete, create a git commit to save the work.

**Process**:

1. **Check git status** to see all changes:
   ```bash
   git status
   git diff --stat
   ```

2. **Review all modified files** from implementation phases

3. **Create commit with descriptive message**:
   - Summarize the iteration goal and what was implemented
   - Reference the iteration number
   - Keep message concise (1-2 sentences)

   Example:
   ```bash
   git add .
   git commit -m "$(cat <<'EOF'
   feat: implement iteration #5 - task validation and comment system

   Adds task validation commands, comment entities, and TUI integration for iteration #5.
   EOF
   )"
   ```

4. **Verify commit was created**:
   ```bash
   git log -1 --oneline
   ```

**Commit Message Guidelines**:
- Start with conventional commit prefix (feat/fix/refactor/docs)
- Include iteration number
- Briefly describe deliverable
- Use HEREDOC format for multi-line messages

**Important**:
- Commit ALL changes from implementation (don't leave uncommitted work)
- Do NOT push to remote (user will push after AC verification)
- Commit should happen even if user will verify AC later

---

### Phase 8: Final Report to User

**IMPORTANT**: Follow the reporting guidelines from CLAUDE.md. Focus on deviations, questions, and issues - not what was implemented (user knows the tasks).

**Report format**:

```markdown
# Implementation Complete: [Iteration Name]

## Status
[One line: Complete / Complete with deviations / Blocked]

## Deviations from Plan
[Only if there were deviations - explain what and why]

## Questions for User
[Only if decisions needed - clear, actionable questions]

## Issues Requiring Attention
[Only if there are blockers or problems]

## Commit
[Commit hash and one-line summary from git log -1 --oneline]
```

**Guidelines**:

**DO**:
- ✅ Highlight deviations from the original plan
- ✅ Emphasize questions that require user decision
- ✅ Call out issues that need user attention
- ✅ Note any blockers or unexpected challenges
- ✅ Use clear, separate blocks for user action items

**DON'T**:
- ❌ Provide detailed summaries of what was implemented (user knows the tasks)
- ❌ Include standard "verify acceptance criteria" instructions (goes without saying)
- ❌ Repeat task descriptions or requirements back to user
- ❌ List every file changed or every test added (user can see git commit)
- ❌ Explain what was supposed to happen (user defined the tasks)

**Example - Good Report**:
```markdown
# Implementation Complete: Iteration #27 TODO Tasks

## Status
Complete with minor deviations

## Deviations from Plan
1. Mocks placed in `application/mocks/` instead of `domain/repositories/mocks/` (AC-482 specifies domain layer)
2. Infrastructure coverage at 52.3% vs 60% target (7.7% gap)

## Questions for User
1. Should mocks stay in application/ or move to domain/repositories/? (affects linter violations)
2. Is 52.3% infrastructure coverage acceptable, or should I add more tests?

## Commit
8b103f8 feat: complete iteration #27 TODO tasks - test architecture refactoring
```

**Example - Bad Report**:
```markdown
# Implementation Complete: Iteration #27 TODO Tasks

## Summary
Successfully implemented all 4 TODO tasks...
[3 paragraphs explaining what was done]

## Implementation Phases
Phase 1: Created mocks...
Phase 2: Refactored tests...
[Detailed breakdown of every step]

## Files Modified
- Created: application/mocks/ (6 files)
- Modified: 10 files
[Long list of every file]

## Next Steps - USER ACTION REQUIRED
1. Review All Acceptance Criteria
2. Verify Each AC
3. Close Tasks
[Standard AC verification instructions]
```

---

## Agent Selection Guide

### Planning Agent
- Always use **general-purpose** (requires design, analysis, planning)

### Implementation Agents
- **junior-dev-executor**: Clear, well-specified work (even if complex)
- **general-purpose**: Exploratory work, research, design decisions

### Verification Agent
- Always use **general-purpose** (requires analysis and review)

---

## Error Handling

**If planning agent fails**:
- Report to user with planning agent's output
- Don't proceed to implementation
- Ask user for guidance

**If implementation phase fails**:
- Don't proceed to dependent phases
- Mark todo as in-progress (not completed)
- Report blocker to user with phase report
- Ask how to proceed

**If verification agent finds issues**:
- Create fix phase
- Delegate to appropriate agent
- Re-run verification
- Only proceed when clean

**If tests or linter fail**:
- Don't report "complete" to user
- Create fix phase
- Resolve all issues before final report

---

## Success Criteria

You succeed when:
- ✅ TODO tasks identified and filtered correctly
- ✅ Planning agent created detailed plan for TODO tasks
- ✅ All implementation phases completed for TODO tasks
- ✅ Verification agent confirms implementation quality
- ✅ All tests pass
- ✅ Zero linter violations
- ✅ TODO tasks marked as "done"
- ✅ Non-TODO tasks (in-progress/review/done/blocked) remain unchanged
- ✅ Git commit created with all changes
- ✅ User notified with clear next steps (verify AC, close iteration)

**Not your responsibility**:
- ❌ Verifying acceptance criteria (user must do this)
- ❌ Closing iteration (user must do this)
- ❌ Deciding if deliverable meets user's needs (user must do this)

---

## Key Principles

1. **TODO tasks only** - Focus exclusively on tasks with status "todo", ignore all others
2. **Orchestrate, don't implement** - Coordinate agents, don't code yourself
3. **Plan first** - Always use planning agent before implementation
4. **Respect dependencies** - Sequential unless planning agent says parallel
5. **Verify thoroughly** - Verification agent checks everything
6. **Track progress** - Update todos and task statuses continuously
7. **Commit always** - Create git commit after implementation is complete
8. **User owns acceptance** - Never verify AC or close iteration yourself
9. **Report clearly** - User must know exactly what to do next

---

Remember: You coordinate a multi-agent workflow. You work ONLY on TODO tasks, leaving in-progress/review/done/blocked tasks unchanged. Planning agent designs, implementation agents execute, verification agent checks. You create a git commit to save all work. You track progress and report results. The user verifies acceptance criteria and closes the iteration.
