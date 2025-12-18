package tools

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

func (r *Registry) registerCrawlingTools() {
	r.Register(&Tool{
		Name:        "katana",
		Category:    "crawling",
		Description: "Next-generation crawling and spidering framework",
		InstallCmd:  "go install github.com/projectdiscovery/katana/cmd/katana@latest",
		ExecuteFunc: executeKatana,
	})

	r.Register(&Tool{
		Name:        "gospider",
		Category:    "crawling",
		Description: "Fast web spider written in Go",
		InstallCmd:  "go install github.com/jaeles-project/gospider@latest",
		ExecuteFunc: executeGospider,
	})

	r.Register(&Tool{
		Name:        "hakrawler",
		Category:    "crawling",
		Description: "Fast web crawler for gathering URLs",
		InstallCmd:  "go install github.com/hakluke/hakrawler@latest",
		ExecuteFunc: executeHakrawler,
	})

	r.Register(&Tool{
		Name:        "waybackurls",
		Category:    "crawling",
		Description: "Fetch URLs from Wayback Machine",
		InstallCmd:  "go install github.com/tomnomnom/waybackurls@latest",
		ExecuteFunc: executeWaybackurls,
	})

	r.Register(&Tool{
		Name:        "gau",
		Category:    "crawling",
		Description: "Fetch known URLs from multiple sources",
		InstallCmd:  "go install github.com/lc/gau/v2/cmd/gau@latest",
		ExecuteFunc: executeGau,
	})

	r.Register(&Tool{
		Name:        "cariddi",
		Category:    "crawling",
		Description: "Fast crawler for URLs and endpoints",
		InstallCmd:  "go install github.com/edoardottt/cariddi/cmd/cariddi@latest",
		ExecuteFunc: executeCariddi,
	})
}

func executeKatana(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "katana.txt")

	cmd := exec.CommandContext(ctx, "katana", "-u", target, "-o", outputFile, "-silent")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("katana failed: %w - %s", err, string(output))
	}

	urls, err := readLines(outputFile)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func executeGospider(ctx context.Context, target, outputDir string) (interface{}, error) {
	cmd := exec.CommandContext(ctx, "gospider", "-s", target, "-o", outputDir, "-q")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("gospider failed: %w - %s", err, string(output))
	}

	return map[string]string{"output_dir": outputDir}, nil
}

func executeHakrawler(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "hakrawler.txt")

	cmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("echo %s | hakrawler > %s", target, outputFile))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("hakrawler failed: %w - %s", err, string(output))
	}

	urls, err := readLines(outputFile)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func executeWaybackurls(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "waybackurls.txt")

	cmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("echo %s | waybackurls > %s", target, outputFile))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("waybackurls failed: %w - %s", err, string(output))
	}

	urls, err := readLines(outputFile)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func executeGau(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "gau.txt")

	cmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("echo %s | gau > %s", target, outputFile))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("gau failed: %w - %s", err, string(output))
	}

	urls, err := readLines(outputFile)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func executeCariddi(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "cariddi.txt")

	cmd := exec.CommandContext(ctx, "cariddi", "-u", target, "-o", outputFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("cariddi failed: %w - %s", err, string(output))
	}

	urls, err := readLines(outputFile)
	if err != nil {
		return nil, err
	}

	return urls, nil
}
