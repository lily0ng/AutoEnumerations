package tools

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

func (r *Registry) registerTechnologyTools() {
	r.Register(&Tool{
		Name:        "webanalyze",
		Category:    "technology",
		Description: "Port of Wappalyzer to automate mass scanning",
		InstallCmd:  "go install github.com/rverton/webanalyze/cmd/webanalyze@latest",
		ExecuteFunc: executeWebanalyze,
	})

	r.Register(&Tool{
		Name:        "whatweb",
		Category:    "technology",
		Description: "Web scanner to identify technologies",
		InstallCmd:  "apt-get install whatweb || brew install whatweb",
		ExecuteFunc: executeWhatweb,
	})

	r.Register(&Tool{
		Name:        "wappalyzer",
		Category:    "technology",
		Description: "Identify technologies on websites",
		InstallCmd:  "npm install -g wappalyzer",
		ExecuteFunc: executeWappalyzer,
	})
}

func executeWebanalyze(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "webanalyze.json")

	cmd := exec.CommandContext(ctx, "webanalyze", "-host", target, "-output", "json", "-worker", "10")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("webanalyze failed: %w", err)
	}

	if err := writeFile(outputFile, output); err != nil {
		return nil, err
	}

	return map[string]string{"output_file": outputFile}, nil
}

func executeWhatweb(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "whatweb.json")

	cmd := exec.CommandContext(ctx, "whatweb", target, "--log-json", outputFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("whatweb failed: %w - %s", err, string(output))
	}

	return map[string]string{"output_file": outputFile}, nil
}

func executeWappalyzer(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "wappalyzer.json")

	cmd := exec.CommandContext(ctx, "wappalyzer", target)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("wappalyzer failed: %w", err)
	}

	if err := writeFile(outputFile, output); err != nil {
		return nil, err
	}

	return map[string]string{"output_file": outputFile}, nil
}
