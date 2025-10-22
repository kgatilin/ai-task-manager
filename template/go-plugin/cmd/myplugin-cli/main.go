package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"example.com/myplugin/internal"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	// Create plugin instance with sample data
	plugin := createPluginWithSampleData()

	// Start event monitoring in background
	go monitorEvents(plugin)

	// Small delay to ensure event monitor starts
	time.Sleep(10 * time.Millisecond)

	// Route command
	command := os.Args[1]
	switch command {
	case "list":
		handleList(plugin)
	case "get":
		handleGet(plugin, os.Args[2:])
	case "update":
		handleUpdate(plugin, os.Args[2:])
	case "info":
		handleInfo(plugin)
	case "types":
		handleTypes(plugin)
	case "help":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}

	// Allow events to be processed before exit
	time.Sleep(50 * time.Millisecond)
}

// createPluginWithSampleData creates a plugin instance with sample items
func createPluginWithSampleData() *internal.ItemPlugin {
	plugin := internal.NewItemPlugin()

	// Initialize plugin
	plugin.SetWorkingDir("/tmp/myplugin-cli")

	// Add sample items
	plugin.AddItem(&internal.Item{
		ID:          "item-1",
		Name:        "Example Item",
		Description: "This is an example item from the external plugin.",
		Tags:        []string{"example", "demo"},
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now().Add(-24 * time.Hour),
	})

	plugin.AddItem(&internal.Item{
		ID:          "item-2",
		Name:        "Another Item",
		Description: "External plugins can run in any language!",
		Tags:        []string{"external", "plugin"},
		CreatedAt:   time.Now().Add(-2 * time.Hour),
		UpdatedAt:   time.Now().Add(-1 * time.Hour),
	})

	plugin.AddItem(&internal.Item{
		ID:          "item-3",
		Name:        "Third Item",
		Description: "Demonstrates multiple entities and querying.",
		Tags:        []string{"query", "demo"},
		CreatedAt:   time.Now().Add(-5 * time.Hour),
		UpdatedAt:   time.Now().Add(-3 * time.Hour),
	})

	// Enable event streaming
	plugin.StartEventStreaming()

	return plugin
}

// monitorEvents captures and displays events emitted by the plugin
func monitorEvents(plugin *internal.ItemPlugin) {
	eventChan := plugin.GetEventChannel()
	for event := range eventChan {
		timestamp := time.Now().Format("15:04:05")
		payloadJSON, _ := json.Marshal(event.Payload)
		fmt.Fprintf(os.Stderr, "[%s] [EVENT] type: %s | source: %s | payload: %s\n",
			timestamp, event.Type, event.Source, string(payloadJSON))
	}
}

// handleList lists all items
func handleList(plugin *internal.ItemPlugin) {
	query := pluginsdk.EntityQuery{
		EntityType: "item",
		Limit:      100,
	}

	items := plugin.QueryEntities(query)

	if len(items) == 0 {
		fmt.Println("No items found.")
		return
	}

	fmt.Printf("Found %d items:\n\n", len(items))
	for _, item := range items {
		prettyPrint(item)
		fmt.Println()
	}
}

// handleGet retrieves a specific item by ID
func handleGet(plugin *internal.ItemPlugin, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: Missing item ID\n")
		fmt.Fprintf(os.Stderr, "Usage: myplugin-cli get <id>\n")
		os.Exit(1)
	}

	itemID := args[0]
	item, err := plugin.GetEntity(itemID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	prettyPrint(item)
}

// handleUpdate updates an item's field
func handleUpdate(plugin *internal.ItemPlugin, args []string) {
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "Error: Missing arguments\n")
		fmt.Fprintf(os.Stderr, "Usage: myplugin-cli update <id> <field> <value>\n")
		fmt.Fprintf(os.Stderr, "Example: myplugin-cli update item-1 name \"New Name\"\n")
		os.Exit(1)
	}

	itemID := args[0]
	field := args[1]
	value := args[2]

	// Build fields map based on field name
	fields := make(map[string]interface{})

	switch field {
	case "name", "description":
		fields[field] = value
	case "tags":
		// Parse comma-separated tags
		tags := strings.Split(value, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
		fields[field] = tags
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown field '%s'. Supported fields: name, description, tags\n", field)
		os.Exit(1)
	}

	updatedItem, err := plugin.UpdateEntity(itemID, fields)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Item updated successfully:")
	prettyPrint(updatedItem)
}

// handleInfo displays plugin information
func handleInfo(plugin *internal.ItemPlugin) {
	info := plugin.GetInfo()

	fmt.Println("Plugin Information:")
	fmt.Printf("  Name:        %s\n", info.Name)
	fmt.Printf("  Version:     %s\n", info.Version)
	fmt.Printf("  Description: %s\n", info.Description)
	fmt.Printf("  Is Core:     %t\n", info.IsCore)

	capabilities := plugin.GetCapabilities()
	fmt.Printf("\nCapabilities:\n")
	for _, cap := range capabilities {
		fmt.Printf("  - %s\n", cap)
	}
}

// handleTypes displays entity types provided by the plugin
func handleTypes(plugin *internal.ItemPlugin) {
	types := plugin.GetEntityTypes()

	if len(types) == 0 {
		fmt.Println("No entity types defined.")
		return
	}

	fmt.Println("Entity Types:")
	for _, entityType := range types {
		fmt.Printf("\n%s %s\n", entityType.Icon, entityType.DisplayName)
		fmt.Printf("  Type:        %s\n", entityType.Type)
		fmt.Printf("  Plural:      %s\n", entityType.DisplayNamePlural)
		fmt.Printf("  Description: %s\n", entityType.Description)
		if len(entityType.Capabilities) > 0 {
			fmt.Printf("  Capabilities: %v\n", entityType.Capabilities)
		}
	}
}

// prettyPrint formats and prints a map as JSON
func prettyPrint(data map[string]interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

// printHelp displays usage information
func printHelp() {
	help := `myplugin-cli - Local testing tool for DarwinFlow plugins

USAGE:
  myplugin-cli <command> [args...]

COMMANDS:
  list                        List all items
  get <id>                    Get item by ID
  update <id> <field> <value> Update item field
  info                        Show plugin information
  types                       Show entity types
  help                        Show this help message

EXAMPLES:
  myplugin-cli list
  myplugin-cli get item-1
  myplugin-cli update item-1 name "New Name"
  myplugin-cli update item-1 description "Updated description"
  myplugin-cli update item-1 tags "tag1,tag2,tag3"
  myplugin-cli info
  myplugin-cli types

NOTES:
  - Events are displayed to stderr with timestamps
  - Command output goes to stdout for easy piping
  - This tool is for local testing only, not for production use
`
	fmt.Fprint(os.Stderr, help)
}
