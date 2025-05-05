package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// these variables will be set by the linker during build using ldflags
var (
	version string = "dev" // default to "dev" if not set
	commit  string = "none"
	date    string = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of nbx",
	Long:  `All software has versions. This is nbx's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Nbx version: %s\n", version)
		fmt.Printf("Git Commit: %s\n", commit)
		fmt.Printf("Build Date: %s\n", date)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
