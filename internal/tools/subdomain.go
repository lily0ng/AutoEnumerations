package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (r *Registry) registerSubdomainTools() {
	r.Register(&Tool{
		Name:        "subfinder",
		Category:    "subdomain",
		Description: "Fast passive subdomain discovery tool",
		InstallCmd:  "go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest",
		ExecuteFunc: executeSubfinder,
	})

	r.Register(&Tool{
		Name:        "assetfinder",
		Category:    "subdomain",
		Description: "Find domains and subdomains related to a given domain",
		InstallCmd:  "go install github.com/tomnomnom/assetfinder@latest",
		ExecuteFunc: executeAssetfinder,
	})

	r.Register(&Tool{
		Name:        "amass",
		Category:    "subdomain",
		Description: "In-depth attack surface mapping and asset discovery",
		InstallCmd:  "go install -v github.com/owasp-amass/amass/v4/...@master",
		ExecuteFunc: executeAmass,
	})

	r.Register(&Tool{
		Name:        "findomain",
		Category:    "subdomain",
		Description: "Fast cross-platform subdomain enumerator",
		InstallCmd:  "wget https://github.com/Findomain/Findomain/releases/latest/download/findomain-linux && chmod +x findomain-linux && mv findomain-linux /usr/local/bin/findomain",
		ExecuteFunc: executeFindomain,
	})

	r.Register(&Tool{
		Name:        "shuffledns",
		Category:    "subdomain",
		Description: "Wrapper around massdns for bruteforce enumeration",
		InstallCmd:  "go install -v github.com/projectdiscovery/shuffledns/cmd/shuffledns@latest",
		ExecuteFunc: executeShuffleDNS,
	})

	r.Register(&Tool{
		Name:        "puredns",
		Category:    "subdomain",
		Description: "Fast domain resolver and subdomain bruteforcing",
		InstallCmd:  "go install github.com/d3mondev/puredns/v2@latest",
		ExecuteFunc: executePureDNS,
	})

	r.Register(&Tool{
		Name:        "crobat",
		Category:    "subdomain",
		Description: "Rapid7 Sonar API client",
		InstallCmd:  "go install github.com/cgboal/sonarsearch/cmd/crobat@latest",
		ExecuteFunc: executeCrobat,
	})
}

func executeSubfinder(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "subfinder.txt")

	cmd := exec.CommandContext(ctx, "subfinder", "-d", target, "-o", outputFile, "-silent")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("subfinder failed: %w - %s", err, string(output))
	}

	subdomains, err := readLines(outputFile)
	if err != nil {
		return nil, err
	}

	return subdomains, nil
}

func executeAssetfinder(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "assetfinder.txt")

	cmd := exec.CommandContext(ctx, "assetfinder", "--subs-only", target)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("assetfinder failed: %w", err)
	}

	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		return nil, err
	}

	subdomains := strings.Split(strings.TrimSpace(string(output)), "\n")
	return subdomains, nil
}

func executeAmass(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "amass.txt")

	cmd := exec.CommandContext(ctx, "amass", "enum", "-passive", "-d", target, "-o", outputFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("amass failed: %w - %s", err, string(output))
	}

	subdomains, err := readLines(outputFile)
	if err != nil {
		return nil, err
	}

	return subdomains, nil
}

func executeFindomain(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "findomain.txt")

	cmd := exec.CommandContext(ctx, "findomain", "-t", target, "-u", outputFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("findomain failed: %w - %s", err, string(output))
	}

	subdomains, err := readLines(outputFile)
	if err != nil {
		return nil, err
	}

	return subdomains, nil
}

func executeShuffleDNS(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "shuffledns.txt")

	cmd := exec.CommandContext(ctx, "shuffledns", "-d", target, "-o", outputFile, "-silent")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("shuffledns failed: %w - %s", err, string(output))
	}

	subdomains, err := readLines(outputFile)
	if err != nil {
		return nil, err
	}

	return subdomains, nil
}

func executePureDNS(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "puredns.txt")

	cmd := exec.CommandContext(ctx, "puredns", "resolve", target, "-w", outputFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("puredns failed: %w - %s", err, string(output))
	}

	subdomains, err := readLines(outputFile)
	if err != nil {
		return nil, err
	}

	return subdomains, nil
}

func executeCrobat(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "crobat.txt")

	cmd := exec.CommandContext(ctx, "crobat", "-s", target)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("crobat failed: %w", err)
	}

	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		return nil, err
	}

	subdomains := strings.Split(strings.TrimSpace(string(output)), "\n")
	return subdomains, nil
}

func readLines(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var result []string
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			result = append(result, line)
		}
	}
	return result, nil
}
