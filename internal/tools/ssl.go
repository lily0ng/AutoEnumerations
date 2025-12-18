package tools

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

func (r *Registry) registerSSLTools() {
	r.Register(&Tool{
		Name:        "tlsx",
		Category:    "ssl",
		Description: "Fast and configurable TLS grabber",
		InstallCmd:  "go install github.com/projectdiscovery/tlsx/cmd/tlsx@latest",
		ExecuteFunc: executeTlsx,
	})

	r.Register(&Tool{
		Name:        "sslyze",
		Category:    "ssl",
		Description: "Fast and powerful SSL/TLS scanning library",
		InstallCmd:  "pip install sslyze",
		ExecuteFunc: executeSslyze,
	})

	r.Register(&Tool{
		Name:        "testssl",
		Category:    "ssl",
		Description: "Testing TLS/SSL encryption",
		InstallCmd:  "git clone https://github.com/drwetter/testssl.sh.git",
		ExecuteFunc: executeTestssl,
	})
}

func executeTlsx(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "tlsx.json")

	cmd := exec.CommandContext(ctx, "tlsx", "-u", target, "-json", "-o", outputFile, "-silent")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("tlsx failed: %w - %s", err, string(output))
	}

	return map[string]string{"output_file": outputFile}, nil
}

func executeSslyze(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "sslyze.json")

	cmd := exec.CommandContext(ctx, "sslyze", "--json_out", outputFile, target)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("sslyze failed: %w - %s", err, string(output))
	}

	return map[string]string{"output_file": outputFile}, nil
}

func executeTestssl(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "testssl.json")

	cmd := exec.CommandContext(ctx, "testssl.sh/testssl.sh", "--jsonfile", outputFile, target)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("testssl failed: %w - %s", err, string(output))
	}

	return map[string]string{"output_file": outputFile}, nil
}
