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
description: Create implementation plan for iteration with architecture-first approach
---

# Iteration Planning Command

You create a **comprehensive implementation plan** for an iteration, applying architecture best practices upfront before any code is written.

**ultrathink throughout** - reason deeply about design choices, layer responsibilities, and how to achieve simplicity while following SOLID/DDD/Clean Architecture principles.

---

## Philosophy: Architecture-First Planning

**Core Principle**: Plan with good architecture in mind from the start, not as an afterthought.

### What Makes a Good Plan?

A good plan:
- **Understands the problem** - Deeply explores requirements and acceptance criteria
- **Explores existing patterns** - Learns from current codebase structure
- **Applies principles proactively** - SOLID, DDD, Clean Architecture from the start
- **Avoids tension** - Designs that are BOTH simple AND principled
- **Provides clear phases** - Breaks work into logical, testable increments
- **Anticipates challenges** - Identifies complexity and dependencies early

### Goal

Create a plan that implementation agents can follow confidently, knowing architecture is sound.

---

## Workflow Overview

1. **Identify target iteration** (current or specified)
2. **Retrieve iteration context** (tasks, ACs, existing documents)
3. **Launch exploration agents** (understand codebase patterns for each task area)
4. **Synthesize architecture approach** (apply best practices to planning)
5. **Create implementation plan** (phases, dependencies, testing)
6. **Save as plan document** (attach to iteration in draft state)
7. **Present to user** (summary and next steps)

---

## Phase 1: Identify Target Iteration

```bash
tm iteration current
```

**If argument provided**: Use specified iteration number
**If no argument**:
- Use current iteration if exists
- Otherwise, find first "planned" iteration by rank:
  ```bash
  tm iteration list | grep planned | head -1
  ```

**If no iterations found**: Exit with message

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
- Task priorities and dependencies (if noted)

**Store**: Task list as `ITERATION_TASKS`

### 2.2: Get All Acceptance Criteria

```bash
tm ac list-iteration $TARGET_ITERATION
```

Parse to understand:
- What acceptance criteria exist for each task
- What functionality is required (descriptions)
- How features will be verified (testing instructions)

**Store**: AC map as `TASK_ACS` (task ID → list of ACs)

### 2.3: Read Existing Documents

```bash
# List documents attached to this iteration
tm doc list --iteration $TARGET_ITERATION

# Read each document (ADRs, existing plans, retrospectives)
tm doc show <doc-id>
```

Parse documents to understand:
- Existing design decisions or constraints
- Previous architectural choices
- Any context or rationale already documented

**Store**: Document summaries

### 2.4: Understand Task Manager Architecture

Read architecture guidelines:

```bash
# Use Read tool to review key architecture docs
# CLAUDE.md - Overall workflow and patterns
# internal/task_manager/CLAUDE.md - Architecture overview
```

Extract key principles:
- Clean Architecture layer rules
- Repository pattern usage
- Use case patterns
- Domain-driven design principles
- Testing expectations (70-80% coverage)

---

## Phase 3: Launch Exploration Agents

**Goal**: Understand current codebase patterns for each task area before planning.

### 3.1: Group Tasks by Area

Analyze `ITERATION_TASKS` and group by:
- Affected layers (domain/application/infrastructure/presentation)
- Affected entities/aggregates
- Feature areas (e.g., task management, iteration workflow, AC verification)

**Store**: Task groups as `TASK_GROUPS` (area → task IDs)

### 3.2: Launch Exploration for Each Area

**For each task group** (can run in parallel):

**Agent**: `Explore` (fast, specialized for codebase exploration)

**Thoroughness level**:
- "quick": 1-2 simple tasks, well-understood area
- "medium": 3-4 tasks or moderate complexity
- "very thorough": 5+ tasks, complex refactoring, or new patterns

**Prompt structure**:

```
Explore the codebase to understand patterns for: [area name]

## Context

Tasks in this area: [task IDs and titles from group]

Acceptance criteria focus:
[List key AC themes for these tasks]

## What to Find

1. **Existing patterns**:
   - How are similar features currently implemented?
   - Which layers are typically involved?
   - What repository methods exist?
   - What application services handle this?

2. **Architecture compliance**:
   - How do existing implementations follow Clean Architecture?
   - Where are domain entities defined?
   - How are dependencies injected?
   - How are use cases structured?

3. **Testing patterns**:
   - How are similar features tested?
   - What mocking strategies are used?
   - What test coverage exists in affected areas?

4. **Potential tension points**:
   - Are there areas where code seems complex for simple problems?
   - Any violations or architectural debt?
   - Any patterns to avoid?

## Search Strategy

Use Glob to find relevant files:
- Domain entities: `internal/task_manager/domain/**/*.go`
- Repositories: `internal/task_manager/infrastructure/**/*repository*.go`
- Use cases: `internal/task_manager/application/**/*.go`
- Commands: `internal/task_manager/presentation/cli/**/*.go`

Use Grep to search for:
- Similar feature implementations
- Entity usage patterns
- Repository method patterns

## Deliverable

Provide exploration report with:
- Existing patterns to follow
- Recommended layer structure
- Potential complexity areas
- Architecture guidance specific to this task area
```

**Execute** (all explorations in parallel if multiple groups):

```
Task tool (one per group):
- subagent_type: "Explore"
- description: "Explore [area] patterns"
- prompt: [above, with task group details]
- model: "haiku" (fast exploration)
```

**Wait for all exploration agents to complete.**

**Store**: Exploration reports as `EXPLORATION_REPORTS` (area → findings)

---

## Phase 4: Synthesize Architecture Approach

**ultrathink**: Review all exploration findings and iteration context to design architecture.

### 4.1: Apply Tension-Driven Architecture Principles

For each task, consider:

**Simplicity Check**:
- Is this a simple problem or complex problem?
- What's the simplest design that could work?
- How do we avoid over-engineering?

**Principle Adherence**:
- Which layer should own this functionality?
- What domain entities/value objects are needed?
- What repository methods are needed?
- What use cases/application services orchestrate?
- Where does presentation/CLI fit in?

**Dependency Flow**:
- Domain imports NOTHING from other layers
- Application imports domain only
- Infrastructure imports domain only (implements interfaces)
- Presentation imports application + domain

**Avoid Tension**:
- Don't fight the architecture - if it feels complex, rethink
- Don't create workarounds - find the right abstraction
- Don't duplicate patterns - follow existing successful patterns
- Don't mix responsibilities - keep layers clean

**Identify Requirement Issues** (CRITICAL for Open Questions section):
- Are requirements clear or ambiguous?
- Do any requirements create unnecessary complexity?
- Are there contradictions or unclear acceptance criteria?
- Could simpler requirements achieve the same goal?
- Do requirements fight against good architecture?

**If tension/ambiguity found**: Document in Open Questions section for user review BEFORE planning implementation.

### 4.2: Design Layer Assignments

For each task:

1. **Domain layer** (pure business logic):
   - New entities or value objects?
   - New domain methods or validation?
   - New repository interfaces?

2. **Application layer** (use cases):
   - New application services?
   - New use case methods?
   - Transaction coordination?

3. **Infrastructure layer** (implementations):
   - New repository implementations?
   - New database queries or migrations?
   - File system operations?

4. **Presentation layer** (CLI):
   - New commands or subcommands?
   - New flags or arguments?
   - Output formatting?

### 4.3: Identify Dependencies and Phases

**Dependencies**:
- Domain entities must exist before repositories
- Repositories must exist before application services
- Application services must exist before commands
- Migrations must run before using new schema

**Phase structure**:
- Group by layer (domain → infrastructure → application → presentation)
- Or group by feature (vertical slice if possible)
- Identify what can be done in parallel
- Identify what must be sequential

### 4.4: Plan Testing Strategy

For each phase/task:
- Unit tests for domain logic (70-80% coverage)
- Integration tests for infrastructure (real DB)
- Use case tests for application (mocked dependencies)
- E2E tests for critical workflows (optional)

---

## Phase 5: Create Implementation Plan Document

### 5.1: Structure Plan Document

**CRITICAL**: Plan focuses on HOW to implement, NOT repeating WHAT needs to be done (already in iteration).

Create implementation strategy document:

```markdown
# Implementation Plan: [Iteration Name]

**Iteration**: #[number] - [name]
**Created**: [date]
**Status**: Draft

---

## Quick Reference

- **Tasks**: [task IDs] (details: `tm task show <id>`)
- **Acceptance Criteria**: View with `tm ac list-iteration [number]`
- **Goal**: [iteration goal from iteration data]
- **Deliverable**: [iteration deliverable from iteration data]

---

## ⚠️ Open Questions

**IMPORTANT**: This section MUST appear first. User reviews this before reading the rest of the plan.

[If no questions, state clearly:]

✅ **No open questions - everything is clear**

All requirements are well-defined and implementation approach is straightforward.

---

[If there ARE questions:]

### Question 1: [Title - e.g., "Clarify task validation requirements"]

**Context**: [What creates ambiguity or tension]

**Current Requirement**: [What the requirement says now]

**Tension/Issue**: [Why this is problematic - creates complexity, unclear, contradictory]

**Options**:
1. **Interpret as [Option A]**
   - Assumption: [what we'd assume]
   - Implementation: [how we'd build it]
   - Risk: [what if assumption is wrong]

2. **Interpret as [Option B]**
   - Assumption: [different interpretation]
   - Implementation: [different approach]
   - Risk: [what if assumption is wrong]

**Recommendation**: [Which interpretation plan suggests and why]

**Decision Needed**: [What user needs to clarify before implementation]

---

### Question 2: [Title - e.g., "Requirement creates architectural tension"]

**Context**: [Describe the requirement]

**Tension**: [Why this requirement fights simplicity or principles]

**Problem**: [What complexity or awkwardness it creates]

**Alternatives**:
1. **Implement as stated**: [consequences - complexity accepted]
2. **Modify requirement to [alternative]**: [how this simplifies, what we lose]

**Recommendation**: [Should we implement as-is or propose requirement change?]

**Decision Needed**: [User chooses: accept complexity OR revise requirement]

---

## Implementation Strategy

### Complexity Assessment

**Overall Complexity**: [Low/Medium/High]

**Rationale**: [Why this complexity rating - what makes it simple or complex?]

**Key Challenges**:
1. [Challenge 1 and why it's challenging]
2. [Challenge 2 and why it's challenging]

### Approach Summary

[2-3 paragraph summary of overall implementation approach]

[Explain the big picture: How will we tackle these tasks? What's the general strategy? Any patterns we'll follow?]

---

## Architecture Design

### Principles Applied

**Clean Architecture**:
[How we'll ensure proper layer separation and dependency flow]

**Domain-Driven Design**:
[What domain concepts/aggregates are involved, how we'll model them]

**Repository Pattern**:
[What repositories we need, what they'll manage]

**Dependency Inversion**:
[How we'll use interfaces and dependency injection]

### Layer Assignments

#### Domain Layer
**New Entities/Value Objects**:
- `EntityName` - [purpose, key responsibilities]
- `ValueObjectName` - [purpose, immutability, validation]

**New Repository Interfaces**:
- `interface EntityRepository` - [key methods: Save, FindByX, etc.]

**Domain Logic**:
- [What business rules/validation belong here]

#### Infrastructure Layer
**Repository Implementations**:
- `sqliteEntityRepository` - [implements EntityRepository]

**Database Changes**:
- Schema migrations: [describe tables/columns to add/modify]
- Indexes: [what needs indexing for performance]
- Queries: [key queries to implement]

**Patterns to Follow**:
[Reference existing repository implementations to follow]

#### Application Layer
**New Services/Use Cases**:
- `EntityService` - [orchestrates what workflows]
- New methods: [list key use case methods]

**Transaction Coordination**:
[How we'll handle multi-step operations, rollback scenarios]

**Patterns to Follow**:
[Reference existing application services to follow]

#### Presentation Layer
**New Commands**:
- `tm entity create` - [flags, arguments, output format]
- `tm entity list` - [filtering, sorting, formatting]

**Modified Commands**:
- `tm existing-command` - [what changes and why]

**Output Formatting**:
[Table/JSON/text format choices, user experience considerations]

### Design Decisions

[For each significant design choice that needs documentation:]

#### Decision: [Title]

**Context**: [What problem requires this decision]

**Options Considered**:
1. **Option A**: [Description]
   - Pros: [advantages]
   - Cons: [disadvantages]
2. **Option B**: [Description]
   - Pros: [advantages]
   - Cons: [disadvantages]

**Decision**: [Chosen option]

**Rationale**: [Why - how it achieves BOTH simplicity AND adherence to principles]

**Trade-offs**: [What we gain, what we lose, what we accept]

**Tension Analysis**: [Does this feel natural or are we fighting the system? If tension exists, why is it acceptable?]

---

## Implementation Phases

[Break work into logical phases - NOT one phase per task, but by implementation strategy]

### Phase 1: [Name - e.g., "Domain Foundation", "Database Schema", "Core Use Cases"]

**Objective**: [What this phase accomplishes - the WHY]

**Scope**: [What gets implemented - the WHAT at high level]
- Implements: [task IDs] (or parts of tasks if tasks span phases)
- Enables: [which ACs become verifiable - by ID]

**Layer Focus**: [Domain/Infrastructure/Application/Presentation/Mixed]

**Implementation Details**:

1. **[Specific work item 1]**
   - Files: [create/modify which files]
   - Pattern: [what pattern to follow, reference to similar code]
   - Considerations: [edge cases, complexity points]

2. **[Specific work item 2]**
   - Files: [create/modify which files]
   - Pattern: [what pattern to follow]
   - Considerations: [edge cases, complexity points]

**Testing**:
- Unit tests: [what scenarios to test]
- Integration tests: [what to test against real DB if infrastructure]
- Coverage target: [percentage]
- Test files: [where tests go]

**Dependencies**:
- Requires: [None / Phase X completed]
- Blocks: [Phase Y, Phase Z]

**Can Run in Parallel With**: [None / Phase X]

**Complexity**: [Low/Medium/High] - [Brief justification]

**Potential Issues**:
- [Issue 1]: [Description and mitigation strategy]
- [Issue 2]: [Description and mitigation strategy]

---

### Phase 2: [Name]

[Repeat structure...]

---

### Phase N: Verification & Polish

**Objective**: Verify implementation meets all requirements and standards

**Scope**:
- Run standard verification workflow (tests, linter, coverage)
- Ensure all ACs are user-verifiable with clear testing instructions
- Update documentation if needed

**Iteration-Specific Checks**:
[Any verification unique to this iteration - e.g.:]
- Test new commands with various flag combinations
- Verify migration runs successfully on empty and populated DBs
- Confirm new UI components render correctly

---

## Testing Strategy

### Iteration-Specific Testing Focus

**Critical Scenarios to Test**:
[What's unique to this iteration that needs thorough testing]
- [Scenario 1 - e.g., "Edge cases in new validation logic"]
- [Scenario 2 - e.g., "Migration with various data states"]
- [Scenario 3 - e.g., "Command flag combinations and error handling"]

**Test Patterns to Follow**:
[Reference existing tests that are good models for this iteration]
- Domain: `path/to/existing_test.go:line` - [why this is a good model]
- Application: `path/to/existing_test.go:line` - [why this is a good model]
- Infrastructure: `path/to/existing_test.go:line` - [why this is a good model]

**Complex Test Cases**:
[Any particularly tricky testing scenarios specific to this iteration]
- [Case 1 and suggested approach]
- [Case 2 and suggested approach]

---

## Risk Analysis

### Risk: [Description of potential problem]

**Likelihood**: [Low/Medium/High]
**Impact**: [Low/Medium/High]
**Phase Affected**: [Phase X]

**Indicators**: [How we'll know if this risk is materializing]

**Mitigation**: [Proactive steps to prevent]

**Contingency**: [What we'll do if it happens anyway]

---

[Repeat for each identified risk...]

---

## References

**Commands to view iteration details**:
- Tasks: `tm task show <task-id>`
- Acceptance Criteria: `tm ac list-iteration <number>`
- Iteration: `tm iteration show <number> --full`

**Architecture Documentation**:
- Overview: `internal/task_manager/CLAUDE.md`
- Workflow: `CLAUDE.md`
- Layer-specific: `internal/task_manager/<layer>/CLAUDE.md`

**Code References** (examples to follow):
[Point to specific files/functions that demonstrate patterns]
```

### 5.2: Use Write Tool to Create Plan File

```bash
# Create filename
PLAN_DATE=$(date +%Y-%m-%d)
PLAN_FILE=".agent/iteration-${TARGET_ITERATION}-plan-${PLAN_DATE}.md"
```

Use Write tool to save plan document to `$PLAN_FILE`

---

## Phase 6: Attach Plan to Iteration

### 6.1: Create Document Record

```bash
tm doc create \
  --title "Implementation Plan: Iteration #${TARGET_ITERATION}" \
  --type plan \
  --from-file "${PLAN_FILE}" \
  --iteration "${TARGET_ITERATION}"
```

**Note**: Document starts in draft state, ready for user review.

### 6.2: Verify Attachment

```bash
tm doc list --iteration $TARGET_ITERATION
```

Confirm plan document appears in list.

---

## Phase 7: Present Summary to User

**DO**:
- ✅ Summarize key architecture decisions
- ✅ Highlight complexity areas or risks
- ✅ Note any open questions requiring user input
- ✅ Present phase structure and dependencies
- ✅ Point to full plan document location
- ✅ Ask if user wants to review or proceed to implementation

**DON'T**:
- ❌ Dump entire plan in chat (it's saved)
- ❌ Make implementation decisions user should review
- ❌ Start implementation without approval

**Format**:

```markdown
# Implementation Plan Created: Iteration #$TARGET_ITERATION

**Plan document**: `.agent/iteration-${TARGET_ITERATION}-plan-${PLAN_DATE}.md`

**Attached to iteration**: Use `tm doc list --iteration $TARGET_ITERATION` to see

## ⚠️ Open Questions

**CRITICAL**: Review this section FIRST.

[If no questions:]
✅ **No open questions - everything is clear**

[If there ARE questions - list each with one-line summary:]
1. **[Question title]** - [Brief description of decision needed]
2. **[Question title]** - [Brief description of decision needed]

→ See plan document for full details on each question.

---

## Plan Summary

**Scope**: [count] tasks ([task IDs only - see `tm task show <id>` for details])
**Strategy**: [count] implementation phases ([sequential/parallel mix])
**Complexity**: [Low/Medium/High]

**Note**: Plan focuses on HOW to implement (architecture, phases, design decisions), not WHAT to implement (tasks/ACs are in iteration data).

## Architecture Approach

**Key Design Decisions**:
1. [Decision 1 - one line summary]
2. [Decision 2 - one line summary]
3. [Decision 3 - one line summary]

**Layer Changes**:
- Domain: [brief summary]
- Infrastructure: [brief summary]
- Application: [brief summary]
- Presentation: [brief summary]

## Phase Structure

[For each phase, one line:]
- Phase 1: [Name] - [Tasks] - [Dependencies]
- Phase 2: [Name] - [Tasks] - [Dependencies]
...

**Parallel opportunities**: [count] phases can run in parallel

## Potential Risks

[Top 2-3 risks if any, otherwise "None identified"]

## Next Steps

Would you like to:
1. Review the full plan (`.agent/iteration-${TARGET_ITERATION}-plan-${PLAN_DATE}.md`)
2. Discuss architecture decisions or make changes
3. Proceed with implementation (`/iteration_implement $TARGET_ITERATION`)
4. Save for later (plan is attached to iteration)
```

---

## Success Criteria

- ✅ Target iteration identified
- ✅ Full iteration context retrieved (tasks, ACs, documents)
- ✅ Exploration agents completed for each task area
- ✅ Architecture approach synthesized (applying best practices)
- ✅ Implementation plan created (phases, dependencies, testing)
- ✅ Plan document saved to `.agent/` directory
- ✅ Plan attached to iteration as document
- ✅ User presented with clear summary and next steps

---

## Key Principles

1. **Architecture-first** - Plan with SOLID/DDD/Clean Architecture from the start
2. **Exploration-based** - Understand existing patterns before planning
3. **Tension-aware** - Seek designs that are BOTH simple AND principled
4. **Layer-conscious** - Clear layer assignments and dependency flow
5. **Phase-structured** - Logical, testable increments with clear dependencies
6. **Risk-aware** - Identify challenges early, plan mitigations
7. **Comprehensive** - Detailed enough for confident implementation
8. **Actionable** - Clear requirements, acceptance criteria, verification steps
9. **Documented** - Captures design rationale and trade-offs
10. **User-centric** - Leaves final decisions to user, presents options clearly

---

Remember: The goal is to create a plan so good that implementation becomes straightforward execution. Architecture thinking happens NOW, not during or after implementation. Use exploration to learn, apply principles to design, create phases that flow naturally. When implementation follows this plan, the code should feel natural - no fighting the architecture, no complex workarounds, just clean implementation of clear design.
