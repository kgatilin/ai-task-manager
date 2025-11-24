# Plugin Metadata

This directory contains the Claude Code plugin metadata.

## plugin.json

Defines the plugin metadata. Skills are auto-discovered from the `skills/` directory.

**Required fields:**
- **name**: Unique identifier for the plugin (string)
- **description**: What the plugin does (string)
- **version**: Semantic version (string, e.g., "1.0.0")

**Optional fields:**
- **author**: Plugin creator (object with `name` field)
- **repository**: Source code location (string URL)
- **homepage**: Plugin homepage (string URL)
- **license**: License type (string, e.g., "MIT")

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
- `author` - Creator (object: `{ "name": "..." }`)
- `homepage` - Plugin website URL (string)
- `repository` - Source code URL (string)
- `license` - License type (string, e.g., "MIT", "Apache-2.0")

### Skills Auto-Discovery

Skills are automatically discovered from the `skills/` directory. Each subdirectory with a `SKILL.md` file is treated as a skill. No explicit declaration needed in plugin.json!

## Adding New Skills

To add a new skill to the plugin:

1. Create directory: `skills/my-new-skill/`
2. Create `skills/my-new-skill/SKILL.md` with proper frontmatter
3. Update version in plugin.json
4. Test the plugin loads correctly

The skill will be automatically discovered! No changes to plugin.json needed beyond version bump.

**Minimal plugin.json:**
```json
{
  "name": "task-manager-plugin",
  "description": "Plugin description",
  "version": "1.0.0",
  "author": {
    "name": "your-name"
  },
  "license": "MIT"
}
```
