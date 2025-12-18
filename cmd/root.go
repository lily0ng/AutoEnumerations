package cmd

import (
	"fmt"
	"os"

	"github.com/autoenumeration/autoenum/internal/config"
	"github.com/autoenumeration/autoenum/internal/engine"
	"github.com/autoenumeration/autoenum/internal/logger"
	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	target     string
	outputDir  string
	configData *config.Config
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "autoenum",
	Short: "AutoEnumeration - Integrated Reconnaissance Framework",
	Long: `AutoEnumeration is a comprehensive reconnaissance framework that integrates
30+ security tools for automated enumeration and vulnerability discovery.

Supports subdomain discovery, port scanning, web enumeration, technology
fingerprinting, and vulnerability scanning with intelligent result correlation.`,
	Version: "1.0.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "target domain or IP")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "./output", "output directory")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func initConfig() {
	var err error
	configData, err = config.LoadConfig(cfgFile)
	if err != nil {
		logger.Error("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	if verbose {
		configData.Verbose = true
	}

	logger.Init(configData.Verbose)

	if outputDir != "" {
		configData.OutputDir = outputDir
	}

	if err := os.MkdirAll(configData.OutputDir, 0755); err != nil {
		logger.Error("Failed to create output directory: %v", err)
		os.Exit(1)
	}
}

func getEngine() *engine.Engine {
	eng, err := engine.NewEngine(configData)
	if err != nil {
		logger.Error("Failed to initialize engine: %v", err)
		os.Exit(1)
	}
	return eng
}

func validateTarget() error {
	if target == "" {
		return fmt.Errorf("target is required (use -t or --target)")
	}
	return nil
}
