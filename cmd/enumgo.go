package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yimincai/nbx/internal/enumgo/generator"
)

// enumgoCmd represents the enumgo command
var enumgoCmd = &cobra.Command{
	Use:   "enumgo",
	Short: "Generate Go enum files from YAML definitions",
	Long: `enumgo is a command-line tool to generate Go enum code from YAML files.

You can specify a single YAML file or a directory containing multiple YAML files.
The generated Go files will be placed in the specified output directory with the given package name.

Example:
  enumgo -in enums.yaml -out ./pkg/enums
  enumgo -in ./yaml/enums -out ./internal/enums -pkg myenums`,

	RunE: func(cmd *cobra.Command, args []string) error { // Use RunE to return an error
		input, _ := cmd.Flags().GetString("in")
		outputDir, _ := cmd.Flags().GetString("out")
		pkgName, _ := cmd.Flags().GetString("pkg") // pkg flag has a default, so no need to check if empty

		// Check if input is a file or directory
		info, err := os.Stat(input)
		if err != nil {
			// Return error directly; cobra handles printing to Stderr and exiting
			return fmt.Errorf("error checking input path '%s': %w", input, err)
		}

		if info.IsDir() {
			// If input is a directory, process all YAML files in the directory
			files, err := os.ReadDir(input)
			if err != nil {
				// Original panicked; return error instead
				return fmt.Errorf("error reading input directory '%s': %w", input, err)
			}

			for _, file := range files {
				// Skip directories within the input directory
				if file.IsDir() {
					continue
				}

				ext := strings.ToLower(filepath.Ext(file.Name()))
				if ext == ".yaml" || ext == ".yml" {
					inputFile := filepath.Join(input, file.Name())
					fmt.Printf("Processing file: %s\n", inputFile) // Optional: print progress
					// Assuming generator.NewEnumFile can return an error
					if err := generator.NewEnumFile(inputFile, outputDir, pkgName); err != nil {
						// If processing one file fails, return the error and stop (or collect errors)
						return fmt.Errorf("error generating enum from '%s': %w", inputFile, err)
					}
				}
			}
		} else {
			// If input is a single file, process the file
			fmt.Printf("Processing file: %s\n", input) // Optional: print progress
			// Assuming generator.NewEnumFile can return an error
			if err := generator.NewEnumFile(input, outputDir, pkgName); err != nil {
				return fmt.Errorf("error generating enum from '%s': %w", input, err)
			}
		}

		fmt.Printf("Enum generation complete. Output files are in %s\n", outputDir) // Optional success message
		return nil                                                                  // Return nil to indicate success
	},
}

func init() {
	// Add the enumgoCmd to the root command (assuming rootCmd is defined elsewhere)
	rootCmd.AddCommand(enumgoCmd)

	// Define flags here
	enumgoCmd.Flags().StringP("in", "i", "", "YAML file or directory containing YAML files defining the enums")
	enumgoCmd.Flags().StringP("out", "o", "", "Directory to output Go files")
	enumgoCmd.Flags().StringP("pkg", "p", "enums", "Target package name") // Default value "enums"

	// Mark required flags so cobra automatically checks them
	enumgoCmd.MarkFlagRequired("in")
	enumgoCmd.MarkFlagRequired("out")

	// Set the command to not print usage on error
	rootCmd.SilenceUsage = true
}
