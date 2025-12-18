package cmd

import (
	"github.com/autoenumeration/autoenum/internal/installer"
	"github.com/autoenumeration/autoenum/internal/logger"
	"github.com/spf13/cobra"
)

var (
	installAll   bool
	installTools []string
	updateTools  bool
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install or update enumeration tools",
	Long: `Install or update the required enumeration tools.
	
This command will download and install the necessary tools for reconnaissance.
Use --all to install all supported tools, or specify individual tools.`,
	RunE: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().BoolVar(&installAll, "all", false, "install all tools")
	installCmd.Flags().StringSliceVar(&installTools, "tools", []string{}, "specific tools to install")
	installCmd.Flags().BoolVar(&updateTools, "update", false, "update existing tools")
}

func runInstall(cmd *cobra.Command, args []string) error {
	logger.Info("Starting tool installation...")

	inst := installer.NewInstaller()

	if installAll {
		logger.Info("Installing all tools...")
		return inst.InstallAll(updateTools)
	}

	if len(installTools) > 0 {
		logger.Info("Installing specified tools: %v", installTools)
		return inst.InstallTools(installTools, updateTools)
	}

	logger.Warning("No tools specified. Use --all or --tools flag")
	return cmd.Help()
}
