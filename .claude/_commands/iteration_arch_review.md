---
allowed-tools:
  - Bash(tm:*)
  - Bash(go-arch-lint:*)
  - Bash(go:*)
  - Bash(git:*)
  - Read
  - Write
  - Glob
  - Grep
  - Task
  - TodoWrite
argument-hint: "[iteration-number]"
description: Review iteration implementation for architectural tension points and design choices
---

# Iteration Architecture Review

You conduct a **design-focused architecture review** of a completed iteration, applying tension-driven architecture philosophy.

**ultrathink throughout** - reason deeply about design choices, tension points, and whether implementation flows naturally or fights the system.

---

## Philosophy: Tension-Driven Architecture

**Core Principle**: Good architecture makes simple things simple while following best practices. Tension signals wrong approach.

### What is Tension?

Tension occurs when:
- Simple problems require complex solutions
- Code needs to be processed/read multiple times
- Workarounds are needed to make things work
- The system "resists" following SOLID/DDD/Clean Architecture principles
- Similar concepts are handled differently in different places
- Execution flow is convoluted or hard to follow

### Key Insight

**Tension is a signal to step back and rethink your approach - not a signal to violate principles OR add complex workarounds.**

If something simple feels complex, the problem isn't "principles are too strict" - the problem is **wrong design/approach**.

### Goal

Find designs that achieve **BOTH** simplicity **AND** adherence to principles. Never compromise on either.

---

## Workflow Overview

1. **Identify target iteration** (current/completed or specified)
2. **Understand implementation** (tasks, ACs, documents, git commits)
3. **Map changes** (files, packages, execution flow)
4. **Check architectural compliance** (SOLID, DDD, Clean Architecture violations)
5. **Identify tension points** (complexity mismatches, flow issues)
6. **Propose alternatives** (designs that are BOTH simple AND principled)
7. **Generate report** (findings + recommendations)
8. **Present to user** (actionable insights)

---

## Phase 1: Identify Target Iteration

```bash
tm iteration current
```

**If argument provided**: Use specified iteration number
**If no argument**:
- Use current iteration if exists
- Otherwise, find most recently completed iteration:
  ```bash
  tm iteration list | grep complete | head -1
  ```

**If no completed iterations found**: Exit with message

**Store as**: `TARGET_ITERATION` (e.g., "38")

---

## Phase 2: Retrieve Iteration Context

### 2.1: Get Iteration Details

```bash
tm iteration show $TARGET_ITERATION --full
```

Parse to understand:
- Iteration name, goal, deliverable
- All tasks (IDs, titles, descriptions, statuses)
- Which tasks were completed in this iteration

### 2.2: Get All Acceptance Criteria

```bash
tm ac list-iteration $TARGET_ITERATION
```

Parse to understand:
- What acceptance criteria exist for each task
- What functionality was required
- Testing instructions (how features are verified)

### 2.3: Read Attached Documents (If Any)

```bash
# List documents attached to this iteration
tm doc list --iteration $TARGET_ITERATION

# Read each document
tm doc show <doc-id>
```

Parse documents to understand:
- Design decisions and rationale
- Architecture patterns intended
- Planning documents and phases
- Any architectural constraints or guidelines

---

## Phase 3: Map Implementation Changes

### 3.1: Find Related Git Commits

```bash
# Search for commits mentioning iteration or tasks
git log --oneline --all --grep="iteration.*$TARGET_ITERATION" | head -20
git log --oneline --all --grep="TM-task-" | head -30

# Get commit details for relevant commits
git show <commit-hash> --stat
```

### 3.2: Identify Changed Files

From commit analysis, identify:
- Which packages were modified (domain/application/infrastructure/presentation)
- Which specific files changed
- What was added vs modified vs deleted

### 3.3: Create Implementation Map

Build mental model:
- **What problem** was being solved (from tasks/ACs)
- **Which packages** were touched
- **What execution flow** was implemented
- **What data flow** was created

Use Glob/Grep to explore:
```bash
# Find files in modified packages
# Use Glob tool with patterns like:
# internal/task_manager/domain/**/*.go
# internal/task_manager/application/**/*.go
```

---

## Phase 4: Launch Deep Analysis Agent

**Agent**: `general-purpose` (requires deep analysis and exploration)

**Prompt structure**:

```
You are conducting a deep architecture review of iteration $TARGET_ITERATION.

## Context

Iteration: $TARGET_ITERATION - [name]
Goal: [iteration goal]
Deliverable: [iteration deliverable]

Tasks implemented: [list of task IDs and titles]

## Your Mission

Analyze the implementation for:
1. **Architectural compliance** (SOLID, DDD, Clean Architecture)
2. **Tension points** (complexity mismatches, awkward flow)
3. **Design quality** (is it naturally simple or fighting the system?)

## Step 1: Read Implementation

For each task in the iteration:

1. Read the task details:
   ```bash
   tm task show <task-id>
   tm ac list <task-id>
   ```

2. Understand what was required (from AC descriptions and testing instructions)

3. Find and read the implementation files using Read tool
   - Focus on files modified for this iteration
   - Understand execution flow
   - Understand data flow

## Step 2: Analyze Execution Flow

For each major feature/change:

**Flow Analysis**:
- How does data flow through the system?
- How does execution/control flow?
- Are there places where flow is convoluted?
- Does anything get processed multiple times?
- Are there workarounds or special cases?

**Complexity Assessment**:
- Is this a simple problem or complex problem?
- Does solution complexity match problem complexity?
- Could a simple problem be solved more simply?

## Step 3: Check Architectural Compliance

```bash
# Run linter
go-arch-lint .

# Get package details if needed
go-arch-lint -format=package internal/task_manager/<layer>

# Get detailed dependency graph if needed
go-arch-lint -detailed -format=markdown .
```

**Check for**:
- Layer violations (domain importing application/infrastructure)
- Dependency inversions not followed
- Interface duplication across packages
- Repository pattern violations
- Use case patterns not followed

**For each violation**:
- Note file:line location
- Understand WHY violation exists
- Is it a tension point? (system fighting principles)

## Step 4: Identify Tension Points

For each feature/change, ask:

**Simplicity Check**:
- Should this be simple? Is it simple?
- If complex: Is complexity justified by problem complexity?
- If simple problem has complex solution: **TENSION POINT**

**Principle Adherence**:
- Does code follow SOLID/DDD/Clean Architecture naturally?
- Or does it fight against principles?
- Are there workarounds to satisfy principles?
- If fighting: **TENSION POINT**

**Consistency Check**:
- Are similar concepts handled similarly?
- Or differently in different places?
- If inconsistent: **TENSION POINT**

**Flow Check**:
- Does execution flow smoothly?
- Or is it convoluted with multiple passes?
- If convoluted: **TENSION POINT**

## Step 5: For Each Tension Point, Propose Alternatives

For each tension point identified:

1. **Describe the tension**: What feels wrong? Where's the awkwardness?
   - Quote specific code with file:line references
   - Explain why it's awkward

2. **Identify root cause**:
   - Wrong abstraction?
   - Wrong requirements?
   - Wrong layer assignment?
   - Wrong pattern choice?

3. **Propose alternative design**:
   - How could we rethink this to be BOTH simple AND principled?
   - Specific code structure or pattern suggestion
   - Show how it eliminates the tension

4. **Trade-off analysis**:
   - Is the tension fundamental (must manage) or eliminable (should redesign)?
   - What do we gain from the alternative?
   - What do we lose?
   - **Recommendation**: Redesign / Accept as-is / Investigate further

## Step 6: Identify Positive Patterns

Note what was done WELL:
- Clean implementations that flow naturally
- Good use of patterns
- Proper layer separation
- Simple solutions to simple problems
- Designs worth replicating

## Step 7: Generate Report

Create comprehensive report (see report structure below).

Include:
- Executive summary with health score
- Architectural compliance (violations with file:line)
- Tension point analysis (each with alternative design proposal)
- Design recommendations (prioritized: Critical/High/Medium/Low)
- Requirements to rethink (if any)
- Positive patterns (celebrate good design)
- Action items (categorized by urgency)

**Return**: Full report content as markdown
```

**Execute**:

```
Task tool:
- subagent_type: "general-purpose"
- description: "Deep architecture review"
- prompt: [above, with iteration details filled in]
- model: "sonnet" (requires deep reasoning)
```

---

## Phase 5: Review Agent Report

When agent returns report:

1. **Review findings**: Understand tension points identified
2. **Validate recommendations**: Do proposed alternatives make sense?
3. **Note critical issues**: Identify must-fix items
4. **Track action items**: Use TodoWrite if specific fixes needed

---

## Phase 6: Create Report Document

Save the report in `.agent/` directory:

```bash
# Create filename with date
REPORT_DATE=$(date +%Y-%m-%d)
REPORT_FILE=".agent/iteration-${TARGET_ITERATION}-arch-review-${REPORT_DATE}.md"
```

Use Write tool to save agent's report to `$REPORT_FILE`

---

## Phase 7: Present Summary to User

**DO**:
- ✅ Highlight critical issues requiring attention
- ✅ Summarize key tension points found
- ✅ Present recommended design alternatives
- ✅ Note any requirements that should be reconsidered
- ✅ Point to full report location
- ✅ Ask if user wants to address issues now or later

**DON'T**:
- ❌ Just dump the full report in chat
- ❌ Overwhelm with details (full report is saved)
- ❌ Make decisions about whether to fix (user decides)

**Format**:

```markdown
# Architecture Review Complete: Iteration $TARGET_ITERATION

**Report saved**: `.agent/iteration-${TARGET_ITERATION}-arch-review-${REPORT_DATE}.md`

## Health Score
[X/10] - [One-line assessment]

## Critical Issues (🔴 Must Address)
[List critical issues - brief, actionable]

## Key Tension Points
[2-3 most significant tension points with one-line summary]

## Design Recommendations
- 🔴 Critical: [count]
- 🟡 Moderate: [count]
- 🟢 Minor: [count]

## Positive Patterns
[1-2 things that were done well]

## Next Steps

Would you like to:
1. Review the full report (`.agent/iteration-${TARGET_ITERATION}-arch-review-${REPORT_DATE}.md`)
2. Address critical issues now
3. Create tasks for fixes
4. Proceed as-is (accept trade-offs)
```

---

## Report Structure (For Analysis Agent)

The analysis agent should generate a report with this structure:

```markdown
# Iteration <NUM> Architecture Review - <DATE>

## Executive Summary

- **Iteration**: #<num> - <name>
- **Review Date**: <date>
- **Health Score**: [1-10]
- **Critical Issues**: <count>
- **Tension Points**: <count>
- **Design Recommendations**: <count>

**Bottom Line**: [One paragraph - is this good architecture or does it need rethinking?]

---

## Iteration Overview

### Tasks Implemented
[List of tasks with IDs and brief descriptions]

### Acceptance Criteria Summary
[Total ACs, how many verified/failed/pending]

### Packages Modified
[List packages/layers touched - domain/application/infrastructure/presentation]

### Key Changes
[High-level summary of what was implemented]

---

## Architectural Compliance

### ✅ Principle Adherence
[What was done well - following SOLID, DDD, Clean Architecture naturally]

Examples:
- Proper dependency inversion (domain defines interfaces, infra implements)
- Good use of repository pattern
- Clean separation of concerns
- Application services orchestrating properly

### ❌ Principle Violations

[For each violation found:]

#### Violation: [Brief description]
- **Type**: [SOLID principle / DDD pattern / Clean Architecture layer]
- **Location**: `file/path.go:line`
- **Issue**: [What's violated and why it matters]
- **Impact**: [Effect on maintainability/testability/complexity]
- **Recommendation**: [How to fix - brief]

### 🔍 Linter Results

```
[Output from go-arch-lint .]
```

**Summary**: [X violations found / Zero violations]

---

## Tension Point Analysis

[For each tension point - typically 2-5 significant ones:]

### Tension Point <N>: <Brief Description>

**Location**: `package/file.go:line-range`

**Observation**:
[Describe what feels wrong/awkward - specific details]

Evidence of tension:
- [ ] Simple problem requiring complex solution?
- [ ] Execution flow convoluted?
- [ ] Multiple processing passes or workarounds?
- [ ] Fighting against architectural principles?
- [ ] Similar concepts handled differently elsewhere?

**Code Example**:
```go
// Quote relevant code section
```

**Root Cause**:
[Why does this tension exist?]
- Wrong abstraction chosen
- Wrong layer assignment
- Wrong pattern for the problem
- Requirements create inherent complexity
- Other: [explain]

**Current Approach**:
[Describe how it's currently implemented and why it creates tension]

**Problem with Current Approach**:
[Why this is not optimal - be specific about the pain points]

**Recommended Alternative Design**:

[Describe alternative approach that achieves BOTH simplicity AND adherence to principles]

```go
// Show proposed structure or pattern
// Be concrete - actual code sketch if possible
```

**How This Eliminates Tension**:
[Explain why the alternative is better - how it makes things simpler AND more principled]

**Trade-off Analysis**:

| Aspect | Current Approach | Proposed Alternative |
|--------|-----------------|---------------------|
| Simplicity | [rating/notes] | [rating/notes] |
| Principle Adherence | [rating/notes] | [rating/notes] |
| Maintainability | [rating/notes] | [rating/notes] |
| Testability | [rating/notes] | [rating/notes] |

**Recommendation**:
- 🔴 **Redesign required** - tension creates significant problems
- 🟡 **Consider refactoring** - tension exists but manageable
- 🟢 **Accept as-is** - tension is fundamental, current approach reasonable

---

## Design Recommendations

### 🔴 Critical - Rethink Required

[Designs that should be reconsidered - create significant complexity or violate core principles]

For each:

#### [Title]
- **Issue**: [What's wrong - specific and actionable]
- **Impact**: [Why this is critical - effect on codebase]
- **Proposed Change**: [Specific design alternative with rationale]
- **Effort Estimate**: [Low/Medium/High - based on scope]
- **Files Affected**: [List specific files]

### 🟡 Moderate - Consider Refactoring

[Designs that work but could be significantly improved]

### 🟢 Minor - Technical Debt

[Small improvements that would be nice but not critical]

---

## Requirements to Rethink

[Are there requirements that, if changed, would eliminate complexity?]

For each:

#### [Requirement Title]

**Current Requirement**:
[What's currently required - specific]

**Tension Created**:
[What complexity/awkwardness this requirement creates]

**Alternative Requirement**:
[How we could change the requirement - be specific]

**Benefits of Alternative**:
- Eliminates: [what complexity goes away]
- Simplifies: [what becomes simpler]
- Maintains: [what core value is preserved]

**Recommendation**:
[Should we propose changing this requirement? Yes/No and why]

---

## Positive Patterns

[Celebrate good designs - what worked well that should be maintained/replicated]

For each:

#### [Pattern Name]
- **Location**: `file/path.go`
- **What's Good**: [Specific positive aspects]
- **Why It Works**: [Why this is a good approach]
- **Replicate Where**: [Other areas that could benefit from this pattern]

Examples:
- Clean separation between domain and infrastructure
- Simple, focused use case implementations
- Effective use of dependency injection
- Good test coverage with clear test structure
- Appropriate abstraction level (not over-engineered)

---

## Action Items

### Immediate (Before Iteration Close)
[Things that should be fixed before considering iteration complete]

- [ ] [Specific action with file references]

### Short-term (Next Sprint/Iteration)
[Refactorings to schedule soon]

- [ ] [Specific action with file references]

### Long-term (Backlog)
[Architectural improvements to consider for future]

- [ ] [Specific action with file references]

### Research Needed
[Areas where more investigation needed before deciding]

- [ ] [What needs investigation and why]

---

## Conclusion

[2-3 paragraphs:]

1. **Overall Assessment**: Is this implementation architecturally sound? Does it flow naturally or fight the system?

2. **Key Takeaway**: What's the most important insight from this review?

3. **Recommendation**:
   - ✅ **Close iteration** - architecture is solid, minor issues acceptable
   - ⚠️ **Fix critical issues first** - address red flags before closing
   - 🔄 **Redesign required** - fundamental issues need rethinking

**Final Health Score**: [X/10] - [Brief justification]
```

---

## Success Criteria

- ✅ Target iteration identified
- ✅ Full iteration context retrieved (tasks, ACs, documents)
- ✅ Implementation changes mapped (commits, files, packages)
- ✅ Deep analysis agent completed review
- ✅ Tension points identified with alternative designs
- ✅ Architectural violations noted with file:line references
- ✅ Design recommendations prioritized
- ✅ Positive patterns celebrated
- ✅ Report saved to `.agent/` directory
- ✅ User presented with actionable summary

---

## Key Principles

1. **Philosophy-driven** - Apply tension-driven architecture thinking
2. **Specific** - File:line references, concrete code examples
3. **Constructive** - Don't just criticize, propose alternatives
4. **Balanced** - Note both problems AND positive patterns
5. **Actionable** - Clear recommendations with priority levels
6. **Trade-off aware** - Analyze options, don't be dogmatic
7. **Flow-focused** - How does data/execution flow? Smooth or convoluted?
8. **Complexity-matched** - Simple problems should have simple solutions
9. **Requirement-questioning** - Sometimes the requirement is the problem
10. **Both/and thinking** - Seek designs that are BOTH simple AND principled

---

Remember: Your goal is NOT to enforce dogmatic adherence to principles. Your goal is to identify where implementation fights against simplicity AND principles, and propose designs that achieve both. Tension is your signal - where you find it, ask: "How could we rethink this to make it flow naturally?"
