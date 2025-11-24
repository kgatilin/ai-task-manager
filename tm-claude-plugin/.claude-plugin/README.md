# Plugin Metadata

This directory contains the Claude Code plugin metadata.

## plugin.json

Defines the plugin structure and metadata:

- **name**: Unique identifier for the plugin (string)
- **description**: What the plugin does (string)
- **version**: Semantic version (string, e.g., "1.0.0")
- **author**: Plugin creator (string)
- **repository**: Source code location (string URL)
- **homepage**: Plugin homepage (string URL)
- **license**: License type (string, e.g., "MIT")
- **skills**: Array of skill definitions with name and path

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
├── README.md                    # User-facing documentation
└── LICENSE                      # MIT License
```

## Components

### Skills

**task-manager-workflows** - Expert guidance for Task Manager CLI workflows
- Auto-invoked when discussing iterations, tasks, acceptance criteria
- Teaches best practices, patterns, and troubleshooting
- Includes comprehensive workflow documentation


## Metadata Fields

### Required
- `name` - Plugin identifier (lowercase, hyphens)
- `description` - What the plugin does
- `version` - Semantic version

### Optional but Recommended
- `author` - Creator name (string)
- `homepage` - Plugin website URL (string)
- `repository` - Source code URL (string)
- `license` - License type (string, e.g., "MIT", "Apache-2.0")

### Skills
- `skills` - Array of skill definitions
  - `name` - Skill identifier (string)
  - `path` - Relative path to skill directory (string)

## Updating plugin.json

When adding new skills:

1. Create the skill directory and SKILL.md file
2. Add entry to `skills` array in plugin.json
3. Update version number
4. Test the plugin loads correctly

Example:
```json
{
  "name": "my-plugin",
  "version": "1.1.0",
  "skills": [
    {
      "name": "existing-skill",
      "path": "skills/existing-skill"
    },
    {
      "name": "my-new-skill",
      "path": "skills/my-new-skill"
    }
  ]
}
```
