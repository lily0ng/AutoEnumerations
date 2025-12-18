package tools

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

func (r *Registry) registerHTTPProbingTools() {
	r.Register(&Tool{
		Name:        "httpx",
		Category:    "http_probing",
		Description: "Fast and multi-purpose HTTP toolkit",
		InstallCmd:  "go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest",
		ExecuteFunc: executeHttpx,
	})

	r.Register(&Tool{
		Name:        "httprobe",
		Category:    "http_probing",
		Description: "Take a list of domains and probe for working HTTP/HTTPS servers",
		InstallCmd:  "go install github.com/tomnomnom/httprobe@latest",
		ExecuteFunc: executeHttpprobe,
	})

	r.Register(&Tool{
		Name:        "feroxbuster",
		Category:    "http_probing",
		Description: "Fast, simple, recursive content discovery tool",
		InstallCmd:  "cargo install feroxbuster",
		ExecuteFunc: executeFeroxbuster,
	})

	r.Register(&Tool{
		Name:        "ffuf",
		Category:    "http_probing",
		Description: "Fast web fuzzer written in Go",
		InstallCmd:  "go install github.com/ffuf/ffuf/v2@latest",
		ExecuteFunc: executeFfuf,
	})

	r.Register(&Tool{
		Name:        "gobuster",
		Category:    "http_probing",
		Description: "Directory/File, DNS and VHost busting tool",
		InstallCmd:  "go install github.com/OJ/gobuster/v3@latest",
		ExecuteFunc: executeGobuster,
	})

	r.Register(&Tool{
		Name:        "dirsearch",
		Category:    "http_probing",
		Description: "Web path scanner",
		InstallCmd:  "git clone https://github.com/maurosoria/dirsearch.git",
		ExecuteFunc: executeDirsearch,
	})
}

func executeHttpx(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "httpx.txt")

	cmd := exec.CommandContext(ctx, "httpx", "-u", target, "-o", outputFile, "-silent", "-tech-detect", "-status-code")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("httpx failed: %w - %s", err, string(output))
	}

	hosts, err := readLines(outputFile)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

func executeHttpprobe(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "httprobe.txt")

	cmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("echo %s | httprobe > %s", target, outputFile))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("httprobe failed: %w - %s", err, string(output))
	}

	hosts, err := readLines(outputFile)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

func executeFeroxbuster(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "feroxbuster.txt")

	cmd := exec.CommandContext(ctx, "feroxbuster", "-u", fmt.Sprintf("http://%s", target), "-o", outputFile, "--silent")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("feroxbuster failed: %w - %s", err, string(output))
	}

	return map[string]string{"output_file": outputFile}, nil
}

func executeFfuf(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "ffuf.json")

	cmd := exec.CommandContext(ctx, "ffuf", "-u", fmt.Sprintf("http://%s/FUZZ", target), "-w", "/usr/share/wordlists/dirb/common.txt", "-o", outputFile, "-of", "json", "-s")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffuf failed: %w - %s", err, string(output))
	}

	return map[string]string{"output_file": outputFile}, nil
}

func executeGobuster(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "gobuster.txt")

	cmd := exec.CommandContext(ctx, "gobuster", "dir", "-u", fmt.Sprintf("http://%s", target), "-w", "/usr/share/wordlists/dirb/common.txt", "-o", outputFile, "-q")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("gobuster failed: %w - %s", err, string(output))
	}

	return map[string]string{"output_file": outputFile}, nil
}

func executeDirsearch(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "dirsearch.txt")

	cmd := exec.CommandContext(ctx, "python3", "dirsearch/dirsearch.py", "-u", fmt.Sprintf("http://%s", target), "-o", outputFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("dirsearch failed: %w - %s", err, string(output))
	}

	return map[string]string{"output_file": outputFile}, nil
}
