package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/autoenumeration/autoenum/internal/logger"
	"github.com/autoenumeration/autoenum/internal/reporter"
	"github.com/spf13/cobra"
)

var (
	scanMode  string
	skipTools []string
	onlyTools []string
	threads   int
	timeout   int
	resume    bool
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run reconnaissance scan on target",
	Long: `Execute a comprehensive reconnaissance scan using configured tools.
	
Scan modes:
  - quick: Fast subdomain and port discovery
  - standard: Full enumeration with web discovery
  - deep: Comprehensive scan with vulnerability detection
  - custom: Use only specified tools`,
	RunE: runScan,
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringVarP(&scanMode, "mode", "m", "standard", "scan mode (quick|standard|deep|custom)")
	scanCmd.Flags().StringSliceVar(&skipTools, "skip", []string{}, "tools to skip")
	scanCmd.Flags().StringSliceVar(&onlyTools, "only", []string{}, "run only these tools")
	scanCmd.Flags().IntVar(&threads, "threads", 10, "number of concurrent threads")
	scanCmd.Flags().IntVar(&timeout, "timeout", 3600, "scan timeout in seconds")
	scanCmd.Flags().BoolVar(&resume, "resume", false, "resume previous scan")
}

func runScan(cmd *cobra.Command, args []string) error {
	if err := validateTarget(); err != nil {
		return err
	}

	logger.Info("Starting AutoEnumeration scan")
	logger.Info("Target: %s", target)
	logger.Info("Mode: %s", scanMode)
	logger.Info("Output: %s", configData.OutputDir)

	eng := getEngine()

	if threads > 0 {
		configData.Concurrency.MaxWorkers = threads
	}

	if timeout > 0 {
		configData.Timeout = timeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Warning("Received interrupt signal, shutting down gracefully...")
		cancel()
	}()

	scanConfig := &ScanConfig{
		Target:    target,
		Mode:      scanMode,
		SkipTools: skipTools,
		OnlyTools: onlyTools,
		Resume:    resume,
	}

	results, err := eng.ExecuteScan(ctx, scanConfig)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	logger.Success("Scan completed successfully")
	logger.Info("Generating reports...")

	rep := reporter.NewReporter(configData.OutputDir)

	if err := rep.GenerateJSON(results); err != nil {
		logger.Error("Failed to generate JSON report: %v", err)
	}

	if err := rep.GenerateHTML(results); err != nil {
		logger.Error("Failed to generate HTML report: %v", err)
	}

	logger.Success("Reports generated in: %s", configData.OutputDir)

	printSummary(results)

	return nil
}

type ScanConfig struct {
	Target    string
	Mode      string
	SkipTools []string
	OnlyTools []string
	Resume    bool
}

func printSummary(results interface{}) {
	logger.Info("\n=== Scan Summary ===")
}
