package tools

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

func (r *Registry) registerCloudTools() {
	r.Register(&Tool{
		Name:        "cloudlist",
		Category:    "cloud",
		Description: "Multi-cloud asset enumeration tool",
		InstallCmd:  "go install github.com/projectdiscovery/cloudlist/cmd/cloudlist@latest",
		ExecuteFunc: executeCloudlist,
	})

	r.Register(&Tool{
		Name:        "s3scanner",
		Category:    "cloud",
		Description: "Scan for open S3 buckets",
		InstallCmd:  "pip install s3scanner",
		ExecuteFunc: executeS3Scanner,
	})

	r.Register(&Tool{
		Name:        "cloud-enum",
		Category:    "cloud",
		Description: "Multi-cloud OSINT tool",
		InstallCmd:  "git clone https://github.com/initstring/cloud_enum.git",
		ExecuteFunc: executeCloudEnum,
	})
}

func executeCloudlist(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "cloudlist.txt")

	cmd := exec.CommandContext(ctx, "cloudlist", "-o", outputFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("cloudlist failed: %w - %s", err, string(output))
	}

	return map[string]string{"output_file": outputFile}, nil
}

func executeS3Scanner(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "s3scanner.txt")

	cmd := exec.CommandContext(ctx, "s3scanner", "scan", "--bucket", target)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("s3scanner failed: %w", err)
	}

	if err := writeFile(outputFile, output); err != nil {
		return nil, err
	}

	return map[string]string{"output_file": outputFile}, nil
}

func executeCloudEnum(ctx context.Context, target, outputDir string) (interface{}, error) {
	outputFile := filepath.Join(outputDir, "cloud-enum.txt")

	cmd := exec.CommandContext(ctx, "python3", "cloud_enum/cloud_enum.py", "-k", target)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("cloud-enum failed: %w", err)
	}

	if err := writeFile(outputFile, output); err != nil {
		return nil, err
	}

	return map[string]string{"output_file": outputFile}, nil
}
