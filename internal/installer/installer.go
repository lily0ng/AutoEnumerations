package installer

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/autoenumeration/autoenum/internal/logger"
	"github.com/autoenumeration/autoenum/internal/tools"
)

type Installer struct {
	registry *tools.Registry
}

func NewInstaller() *Installer {
	return &Installer{
		registry: tools.GetRegistry(),
	}
}

func (i *Installer) InstallAll(update bool) error {
	allTools := i.registry.GetAllTools()

	logger.Info("Installing %d tools...", len(allTools))

	failed := 0
	for _, tool := range allTools {
		if err := i.installTool(tool, update); err != nil {
			logger.Error("Failed to install %s: %v", tool.Name, err)
			failed++
		}
	}

	if failed > 0 {
		return fmt.Errorf("%d tools failed to install", failed)
	}

	logger.Success("All tools installed successfully!")
	return nil
}

func (i *Installer) InstallTools(toolNames []string, update bool) error {
	failed := 0

	for _, name := range toolNames {
		tool := i.registry.GetTool(name)
		if tool == nil {
			logger.Warning("Tool not found: %s", name)
			continue
		}

		if err := i.installTool(tool, update); err != nil {
			logger.Error("Failed to install %s: %v", name, err)
			failed++
		}
	}

	if failed > 0 {
		return fmt.Errorf("%d tools failed to install", failed)
	}

	logger.Success("Tools installed successfully!")
	return nil
}

func (i *Installer) installTool(tool *tools.Tool, update bool) error {
	if tool.IsInstalled() && !update {
		logger.Info("✓ %s is already installed", tool.Name)
		return nil
	}

	logger.Info("Installing %s...", tool.Name)

	if tool.InstallCmd == "" {
		return fmt.Errorf("no install command defined for %s", tool.Name)
	}

	parts := strings.Fields(tool.InstallCmd)
	if len(parts) == 0 {
		return fmt.Errorf("invalid install command for %s", tool.Name)
	}

	var cmd *exec.Cmd
	if strings.Contains(tool.InstallCmd, "&&") || strings.Contains(tool.InstallCmd, "||") {
		cmd = exec.Command("sh", "-c", tool.InstallCmd)
	} else {
		cmd = exec.Command(parts[0], parts[1:]...)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("installation failed: %w - %s", err, string(output))
	}

	logger.Success("✓ %s installed successfully", tool.Name)
	return nil
}

func (i *Installer) CheckInstalled() map[string]bool {
	allTools := i.registry.GetAllTools()
	status := make(map[string]bool)

	for _, tool := range allTools {
		status[tool.Name] = tool.IsInstalled()
	}

	return status
}
