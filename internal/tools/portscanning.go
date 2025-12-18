package tools

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

func (r *Registry) registerPortScanningTools() {
	r.Register(&Tool{
		Name:        "naabu",
		Category:    "port_scanning",
		Description: "Fast port scanner written in Go",
		InstallCmd:  "go install -v github.com/projectdiscovery/naabu/v2/cmd/naabu@latest",
		ExecuteFunc: executeNaabu,
	})

	r.Register(&Tool{
		Name:        "nmap",
		Category:    "port_scanning",
		Description: "Network exploration tool and security scanner",
		InstallCmd:  "apt-get install nmap || brew install nmap",
		ExecuteFunc: executeNmap,
	})

	r.Register(&Tool{
		Name:        "rustscan",
		Category:    "port_scanning",
		Description: "Modern port scanner",
		InstallCmd:  "cargo install rustscan",
		ExecuteFunc: executeRustScan,
	})

	r.Register(&Tool{
		Name:        "masscan",
		Category:    "port_scanning",
		Description: "TCP port scanner, spews SYN packets asynchronously",
		InstallCmd:  "apt-get install masscan || brew install masscan",
		ExecuteFunc: executeMasscan,
	})
}

func executeNaabu(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "naabu.txt")

	cmd := exec.CommandContext(ctx, "naabu", "-host", target, "-o", outputFile, "-silent", "-top-ports", "1000")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("naabu failed: %w - %s", err, string(output))
	}

	ports, err := readLines(outputFile)
	if err != nil {
		return nil, err
	}

	return ports, nil
}

func executeNmap(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "nmap.xml")

	cmd := exec.CommandContext(ctx, "nmap", "-sV", "-T4", "-oX", outputFile, target)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("nmap failed: %w - %s", err, string(output))
	}

	return map[string]string{"output_file": outputFile}, nil
}

func executeRustScan(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "rustscan.txt")

	cmd := exec.CommandContext(ctx, "rustscan", "-a", target, "--ulimit", "5000")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("rustscan failed: %w", err)
	}

	if err := writeFile(outputFile, output); err != nil {
		return nil, err
	}

	return map[string]string{"output_file": outputFile}, nil
}

func executeMasscan(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "masscan.txt")

	cmd := exec.CommandContext(ctx, "masscan", target, "-p1-65535", "--rate=10000", "-oL", outputFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("masscan failed: %w - %s", err, string(output))
	}

	ports, err := readLines(outputFile)
	if err != nil {
		return nil, err
	}

	return ports, nil
}
