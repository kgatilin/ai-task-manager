# Claude Code Plugin Marketplace

This directory contains the marketplace configuration for plugins available in this repository.

## Available Plugins

### tm-claude-plugin
Master Task Manager workflows with expert AI guidance for iterations, tasks, and acceptance criteria.

**Location**: `tm-claude-plugin/`

## Installing Plugins

Users can install plugins from this repository using the `/plugin` command in Claude Code:

```bash
# Add this repository as a marketplace
/plugin marketplace add kgatilin/ai-task-manager

# Install the Task Manager plugin
/plugin install tm-claude-plugin@ai-task-manager
```

## Marketplace Structure

```
.claude-plugin/
├── marketplace.json    # Marketplace metadata and plugin catalog
└── README.md          # This file

tm-claude-plugin/      # The actual plugin
├── .claude-plugin/
│   └── plugin.json
├── skills/
├── commands/
└── ...
```

## marketplace.json Format

The `marketplace.json` file defines available plugins:

```json
{
  "name": "marketplace-name",
  "plugins": [
    {
      "name": "plugin-name",
      "path": "relative/path/to/plugin",
      "description": "Plugin description",
      "version": "1.0.0"
    }
  ]
}
```

Each plugin's `path` should point to a directory containing `.claude-plugin/plugin.json`.
