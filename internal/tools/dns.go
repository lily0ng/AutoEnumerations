package tools

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

func (r *Registry) registerDNSTools() {
	r.Register(&Tool{
		Name:        "dnsx",
		Category:    "dns",
		Description: "Fast and multi-purpose DNS toolkit",
		InstallCmd:  "go install -v github.com/projectdiscovery/dnsx/cmd/dnsx@latest",
		ExecuteFunc: executeDnsx,
	})

	r.Register(&Tool{
		Name:        "dnsrecon",
		Category:    "dns",
		Description: "DNS enumeration script",
		InstallCmd:  "apt-get install dnsrecon || pip install dnsrecon",
		ExecuteFunc: executeDnsrecon,
	})

	r.Register(&Tool{
		Name:        "massdns",
		Category:    "dns",
		Description: "High-performance DNS stub resolver",
		InstallCmd:  "git clone https://github.com/blechschmidt/massdns.git && cd massdns && make",
		ExecuteFunc: executeMassdns,
	})

	r.Register(&Tool{
		Name:        "zdns",
		Category:    "dns",
		Description: "Fast CLI DNS lookup tool",
		InstallCmd:  "go install github.com/zmap/zdns@latest",
		ExecuteFunc: executeZdns,
	})
}

func executeDnsx(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "dnsx.txt")

	cmd := exec.CommandContext(ctx, "dnsx", "-d", target, "-o", outputFile, "-silent", "-resp")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("dnsx failed: %w - %s", err, string(output))
	}

	records, err := readLines(outputFile)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func executeDnsrecon(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "dnsrecon.json")

	cmd := exec.CommandContext(ctx, "dnsrecon", "-d", target, "-j", outputFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("dnsrecon failed: %w - %s", err, string(output))
	}

	return map[string]string{"output_file": outputFile}, nil
}

func executeMassdns(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "massdns.txt")

	cmd := exec.CommandContext(ctx, "massdns", "-r", "/etc/resolv.conf", "-t", "A", "-o", "S", "-w", outputFile, target)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("massdns failed: %w - %s", err, string(output))
	}

	return map[string]string{"output_file": outputFile}, nil
}

func executeZdns(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "zdns.json")

	cmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("echo %s | zdns A > %s", target, outputFile))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("zdns failed: %w - %s", err, string(output))
	}

	return map[string]string{"output_file": outputFile}, nil
}
