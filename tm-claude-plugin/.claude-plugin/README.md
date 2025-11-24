# Plugin Metadata

This directory contains the Claude Code plugin metadata.

## plugin.json

Defines the plugin structure, components, and metadata:

- **name**: Unique identifier for the plugin
- **description**: What the plugin does (shown in plugin listings)
- **version**: Semantic version (1.0.0)
- **author**: Plugin creator information
- **repository**: Source code location
- **components**: Skills and commands bundled in this plugin

## Plugin Structure

```
tm-claude-plugin/
├── .claude-plugin/
│   ├── plugin.json              # THIS FILE - Plugin metadata
│   └── README.md                # This documentation
├── skills/                      # Skills (auto-invoked by Claude)
│   └── task-manager-workflows/
│       ├── SKILL.md            # Main skill definition
│       └── references/         # Detailed documentation
├── commands/                    # Slash commands (explicit invocation)
│   ├── tm-current.md
│   ├── tm-status.md
│   └── tm-validate.md
├── README.md                    # User-facing documentation
├── INSTALLATION.md              # Installation guide
└── LICENSE                      # MIT License
```

## Components

### Skills

**task-manager-workflows** - Expert guidance for Task Manager CLI workflows
- Auto-invoked when discussing iterations, tasks, acceptance criteria
- Teaches best practices, patterns, and troubleshooting
- Includes comprehensive workflow documentation

### Commands

**tm-current** - Show current or next planned iteration
**tm-validate** - Validate acceptance criteria in iteration
**tm-status** - Comprehensive project status overview

## Metadata Fields

### Required
- `name` - Plugin identifier (lowercase, hyphens)
- `description` - What the plugin does
- `version` - Semantic version

### Optional but Recommended
- `author` - Creator information
- `homepage` - Plugin website/documentation
- `repository` - Source code location
- `license` - License type (MIT, Apache, etc.)
- `keywords` - Search/discovery terms

### Components
- `skills` - Array of skill definitions
  - `name` - Skill identifier
  - `path` - Relative path to skill directory
- `commands` - Array of slash command definitions
  - `name` - Command name (invoked as /name)
  - `path` - Relative path to command markdown

## Updating plugin.json

When adding new skills or commands:

1. Create the skill/command file
2. Add entry to `components` in plugin.json
3. Update version number
4. Test the plugin loads correctly

Example:
```json
{
  "components": {
    "skills": [
      {
        "name": "my-new-skill",
        "path": "skills/my-new-skill"
      }
    ],
    "commands": [
      {
        "name": "my-command",
        "path": "commands/my-command.md"
      }
    ]
  }
}
```
