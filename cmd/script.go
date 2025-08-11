/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/byte4cat/nbx/internal/script"
	"github.com/byte4cat/nbx/pkg/clog"
	"github.com/spf13/cobra"
)

// scriptsCmd represents the scripts command
var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Script",
	Long: `Script is a collection of scripts that can be used to automate tasks.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// 取得 flags
		platform, _ := cmd.Flags().GetString("platform")
		scriptName, _ := cmd.Flags().GetString("script")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// 查找平台
		platformScripts, ok := script.Registry[platform]
		if !ok {
			clog.Error("Unsupported platform: %s", platform)
			clog.Info("Supported platforms")
			for p := range script.Registry {
				clog.Item(p)
			}
			clog.Info("Hint: try `--platform <platform> --script <script>`")
			return
		}

		// 查找腳本
		fn, ok := platformScripts[scriptName]
		if !ok {
			clog.Error("Unknown script: %s for platform: %s", scriptName, platform)
			clog.Info("Available scripts for platform '%s':", platform)
			for name := range platformScripts {
				clog.Item(name)
			}
			clog.Info("Hint: try `--platform %s --script <name>`", platform)
			return
		}

		if dryRun {
			clog.Info("Dry run script: %s on platform: %s", scriptName, platform)
			clog.Console(fn.DryRun())
		} else {
			clog.Info("Running script: %s on platform: %s", scriptName, platform)
			fn.Run()
		}
	},
}

func init() {
	rootCmd.AddCommand(scriptCmd)

	scriptCmd.Flags().StringP("platform", "p", "", "The platform to run the script on")
	scriptCmd.Flags().StringP("script", "s", "", "The script to run")
	scriptCmd.Flags().Bool("dry-run", false, "Print the script without executing")

	scriptCmd.MarkFlagRequired("platform")
}
