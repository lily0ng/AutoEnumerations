package cmd

import (
	"fmt"

	"github.com/autoenumeration/autoenum/internal/tools"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listCategory string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available tools and categories",
	Long:  `Display all available enumeration tools organized by category.`,
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listCategory, "category", "c", "", "filter by category")
}

func runList(cmd *cobra.Command, args []string) error {
	registry := tools.GetRegistry()

	categories := registry.GetCategories()

	if listCategory != "" {
		return listToolsByCategory(registry, listCategory)
	}

	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan)

	fmt.Println("\n=== AutoEnumeration Tools ===\n")

	for _, cat := range categories {
		green.Printf("📦 %s\n", cat)
		toolList := registry.GetToolsByCategory(cat)
		for _, tool := range toolList {
			status := "❌"
			if tool.IsInstalled() {
				status = "✅"
			}
			cyan.Printf("  %s %s - %s\n", status, tool.Name, tool.Description)
		}
		fmt.Println()
	}

	return nil
}

func listToolsByCategory(registry *tools.Registry, category string) error {
	toolList := registry.GetToolsByCategory(category)

	if len(toolList) == 0 {
		return fmt.Errorf("category not found: %s", category)
	}

	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan)

	green.Printf("\n=== %s Tools ===\n\n", category)

	for _, tool := range toolList {
		status := "❌"
		if tool.IsInstalled() {
			status = "✅"
		}
		cyan.Printf("%s %s\n", status, tool.Name)
		fmt.Printf("  Description: %s\n", tool.Description)
		fmt.Printf("  Install: %s\n", tool.InstallCmd)
		fmt.Println()
	}

	return nil
}
