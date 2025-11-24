# Task Manager Claude Code Plugin

**Master Task Manager workflows with expert AI guidance**

This Claude Code skill provides comprehensive guidance for using the Task Manager CLI (`tm`) - a hierarchical roadmap and task management tool built with Clean Architecture. It teaches Claude how to help you plan, execute, and verify software development work using iterations, tasks, and acceptance criteria.

## What This Plugin Provides

### 🎯 Expert Skill: Task Manager Workflows

This plugin provides a comprehensive skill that teaches Claude how to work with Task Manager. Claude learns:
- When to use `tm iteration current`, `tm task list`, `tm ac add`, and other commands
- How to write good acceptance criteria (end-user verifiable, not implementation details)
- Task granularity guidelines (user-facing features vs implementation details)
- Iteration workflow patterns (planning → execution → validation)
- Document management (ADRs, plans, retrospectives)
- Common troubleshooting scenarios

### 📚 Reference Documentation

Complete command reference for all `tm` operations.

## Installation

Use Claude Code's built-in plugin manager:

```bash
# In Claude Code (from any project):

# 1. Add the marketplace
/plugin marketplace add kgatilin/ai-task-manager

# 2. Install the plugin
/plugin install tm-claude-plugin@ai-task-manager
```

**That's it!** The plugin is now available in your project.

### Plugin Management

```bash
/plugin                                         # Open interactive plugin manager
/plugin disable tm-claude-plugin                # Temporarily disable
/plugin enable tm-claude-plugin                 # Re-enable
/plugin uninstall tm-claude-plugin@ai-task-manager  # Remove completely
```

## Verifying Installation

After installation, check that the plugin is loaded:

1. In Claude Code, ask: "What Task Manager workflows do you know?"
2. Claude should mention iteration workflows, acceptance criteria patterns, etc.

## Usage

### Using the Skill (Automatic)

Claude automatically invokes the skill when you discuss Task Manager workflows:

**Examples:**
```
You: "How do I start working on an iteration?"
Claude: [Explains tm iteration current, tm iteration show, etc.]

You: "Help me create good acceptance criteria for this task"
Claude: [Guides you through AC best practices with examples]

You: "What should I work on next?"
Claude: [Runs tm iteration current and explains priority]
```

## What You'll Learn

### Core Workflows

1. **Understanding What to Work On**
   - Always start with `tm iteration current`
   - View iteration details and task breakdown
   - Check acceptance criteria before starting

2. **Task Lifecycle Management**
   - Status transitions: `todo` → `in-progress` → `done`
   - When to update status (do it as you work, not in batches)
   - Agent vs user responsibilities

3. **Acceptance Criteria Best Practices**
   - Write end-user verifiable behavior, not implementation details
   - Separate `--description` (WHAT) from `--testing-instructions` (HOW)
   - Focus on observable outcomes
   - Include exact commands for verification

4. **Iteration Workflow**
   - Planning iterations with clear goals and deliverables
   - Adding tasks to iterations
   - Validation process (`tm iteration validate`)
   - Handling failed acceptance criteria

5. **Document Management**
   - Creating ADRs for architectural decisions
   - Attaching documents to tracks or iterations
   - Using plans and retrospectives

### Task Granularity Guidelines

Learn when tasks are too granular:
- ✅ Good: "Add task comment commands" (complete user feature)
- ❌ Too granular: "Create TaskComment entity" (implementation detail)

If you can't write end-user verifiable AC, the task needs to be merged into a larger feature.

### Architecture Context

Understand Clean Architecture layers:
- **Domain**: Pure business logic
- **Application**: Use cases and orchestration
- **Infrastructure**: Technical implementations (database, file system)
- **Presentation**: CLI and TUI

This helps you focus AC on user-facing features, not internal architecture.

## Prerequisites

This plugin requires Task Manager (`tm`) CLI to be installed.

### Installing Task Manager

**Quick install (requires Go 1.21+):**

```bash
go install github.com/kgatilin/ai-task-manager/cmd/tm@latest
```

This installs the `tm` command to your `$GOPATH/bin` (usually `~/go/bin`).

**Verify installation:**

```bash
tm --version
# Output: tm version 1.0.0
```

**Troubleshooting:**

If `tm` command is not found after installation, ensure `$GOPATH/bin` is in your PATH:

```bash
# Add to ~/.bashrc, ~/.zshrc, or equivalent
export PATH="$PATH:$(go env GOPATH)/bin"
```

**Alternative installations:**

```bash
# Build from source
git clone https://github.com/kgatilin/ai-task-manager.git
cd ai-task-manager
make build
```

See the [Task Manager README](https://github.com/kgatilin/ai-task-manager) for more installation options including Docker containers with headless binary.

## Example Workflows

### Starting New Work
```
You: "Show me what I should work on"
Claude: [Runs tm iteration current and explains tasks]

You: "I want to start task TM-task-5"
Claude: [Runs tm task update TM-task-5 --status in-progress and shows AC]
```

### Creating Good Acceptance Criteria
```
You: "Help me add acceptance criteria for implementing user authentication"
Claude: [Guides you through creating user-verifiable AC with proper structure]
```

## Extending the Skill

Add custom patterns to the skill:

```bash
# Add to references
echo "# Custom Patterns" > .claude/plugins/tm-claude-plugin/skills/task-manager-workflows/references/custom-patterns.md

# Reference from SKILL.md
vim .claude/plugins/tm-claude-plugin/skills/task-manager-workflows/SKILL.md
```

## Troubleshooting

### Plugin Not Loading

Check plugin location:
```bash
# Project-level
ls .claude/plugins/tm-claude-plugin/.claude-plugin/plugin.json

# Global
ls ~/.claude/plugins/tm-claude-plugin/.claude-plugin/plugin.json
```

Verify `plugin.json` is valid JSON:
```bash
cat .claude/plugins/tm-claude-plugin/.claude-plugin/plugin.json | python -m json.tool
```

### Skill Not Invoked

Check the skill description in `SKILL.md` frontmatter—it should clearly describe when to use the skill:

```yaml
---
name: task-manager-workflows
description: Expert guidance for Task Manager CLI (tm). Helps navigate iterations, tasks, acceptance criteria...
---
```

The description is how Claude decides when to invoke your skill.

## Contributing

Improvements welcome! To contribute:

1. Fork the [ai-task-manager repository](https://github.com/kgatilin/ai-task-manager)
2. Make changes in `tm-claude-plugin/`
3. Test locally by copying to `.claude/plugins/`
4. Submit a pull request

## Support

- **Task Manager Issues**: https://github.com/kgatilin/ai-task-manager/issues
- **Plugin Issues**: https://github.com/kgatilin/ai-task-manager/issues (tag: "plugin")
- **Claude Code Docs**: https://code.claude.com/docs

## License

MIT License - See LICENSE file in repository root.

## Related Resources

- **Task Manager Repository**: https://github.com/kgatilin/ai-task-manager
- **Claude Code Skills**: https://code.claude.com/docs/en/skills
- **Claude Code Plugins**: https://code.claude.com/docs/en/plugins
- **Slash Commands**: https://code.claude.com/docs/en/slash-commands

---

**Happy task managing! 🚀**
