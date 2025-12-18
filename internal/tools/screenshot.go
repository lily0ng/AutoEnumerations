package tools

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

func (r *Registry) registerScreenshotTools() {
	r.Register(&Tool{
		Name:        "gowitness",
		Category:    "screenshot",
		Description: "Web screenshot utility using Chrome Headless",
		InstallCmd:  "go install github.com/sensepost/gowitness@latest",
		ExecuteFunc: executeGowitness,
	})

	r.Register(&Tool{
		Name:        "aquatone",
		Category:    "screenshot",
		Description: "Tool for visual inspection of websites",
		InstallCmd:  "go install github.com/michenriksen/aquatone@latest",
		ExecuteFunc: executeAquatone,
	})

	r.Register(&Tool{
		Name:        "eyewitness",
		Category:    "screenshot",
		Description: "Take screenshots of websites",
		InstallCmd:  "git clone https://github.com/FortyNorthSecurity/EyeWitness.git",
		ExecuteFunc: executeEyewitness,
	})
}

func executeGowitness(ctx context.Context, target, outputDir string) (interface{}, error) {
	screenshotDir := filepath.Join(outputDir, "screenshots")

	cmd := exec.CommandContext(ctx, "gowitness", "single", target, "--screenshot-path", screenshotDir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("gowitness failed: %w - %s", err, string(output))
	}

	return map[string]string{"screenshot_dir": screenshotDir}, nil
}

func executeAquatone(ctx context.Context, target, outputDir string) (interface{}, error) {
	screenshotDir := filepath.Join(outputDir, "aquatone")

	cmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("echo %s | aquatone -out %s", target, screenshotDir))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("aquatone failed: %w - %s", err, string(output))
	}

	return map[string]string{"screenshot_dir": screenshotDir}, nil
}

func executeEyewitness(ctx context.Context, target, outputDir string) (interface{}, error) {
	screenshotDir := filepath.Join(outputDir, "eyewitness")

	cmd := exec.CommandContext(ctx, "python3", "EyeWitness/Python/EyeWitness.py", "--single", target, "-d", screenshotDir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("eyewitness failed: %w - %s", err, string(output))
	}

	return map[string]string{"screenshot_dir": screenshotDir}, nil
}
