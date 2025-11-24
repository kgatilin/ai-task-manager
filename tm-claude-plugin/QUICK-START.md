# Quick Start

## Prerequisites

Install Task Manager CLI:
```bash
go install github.com/kgatilin/ai-task-manager/cmd/tm@latest
```

Verify: `tm --version` should show `tm version 1.0.0`

## Install the Plugin

```bash
/plugin marketplace add kgatilin/ai-task-manager
/plugin install tm-claude-plugin@ai-task-manager
```

## Use It

Just talk to Claude about Task Manager:

```
"How do I start an iteration?"
"Help me write acceptance criteria"
"What should I work on next?"
```

Claude will automatically use the Task Manager skill knowledge to help you.

## Manage It

```bash
/plugin                    # Browse and manage plugins
/plugin disable tm-claude-plugin
/plugin enable tm-claude-plugin
/plugin uninstall tm-claude-plugin@ai-task-manager
```

That's it! 🚀
