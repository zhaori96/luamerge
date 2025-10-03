package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"luamerge/internal/config"
	"luamerge/internal/merger"
	"luamerge/internal/preservation"
	tmpl "luamerge/internal/template"

	"github.com/spf13/cobra"
)

var (
	inputDir string
	version  = "dev"
)

var rootCmd = &cobra.Command{
	Use:     "luamerge",
	Short:   "Merge Lua table files based on settings.json configuration",
	Version: version,
	Long: `luamerge is a tool to merge Lua table files using job-based configuration.
It reads settings.json from the input directory and processes all configured jobs.
Each job can merge multiple tables from a pair of Lua files.`,
	Run: func(cmd *cobra.Command, args []string) {
		inputPath, _ := cmd.Flags().GetString("inputs")

		// Load configuration from settings.json
		settings, err := config.LoadSettingsFromInput(inputPath)
		if err != nil {
			log.Fatalf("❌ Error loading settings.json: %v", err)
		}

		// Load embedded template
		tpl, err := template.New("lua").Parse(tmpl.LuaTemplate)
		if err != nil {
			log.Fatalf("❌ Error loading embedded template: %v", err)
		}

		fmt.Printf("🚀 luamerge - Processing %d job(s)...\n\n", len(settings.Jobs))

		// Process each job
		for i, job := range settings.Jobs {
			jobName := job.Name
			if jobName == "" {
				jobName = fmt.Sprintf("Job %d", i+1)
			}

			fmt.Printf("[%d/%d] %s\n", i+1, len(settings.Jobs), jobName)

			// Resolve job paths
			basePath, sourcePath, outputPath, err := config.ResolveJobPaths(job, inputPath)
			if err != nil {
				log.Fatalf("❌ Error resolving paths for job '%s': %v", jobName, err)
			}

			// Check if unmerged items should be preserved
			keepUnmerged := job.GetKeepUnmergedItems(settings.Options)

			// Create output directory if it doesn't exist
			outputDir := filepath.Dir(outputPath)
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				log.Fatalf("❌ Error creating output directory '%s': %v", outputDir, err)
			}

			var outputContent string

			// Normalize tables configuration
			tablesConfig := job.GetTablesConfig()

			if keepUnmerged {
				// Mode: Preserve original file and replace only merged tables
				fmt.Printf("  ℹ️  Mode: Preserving unspecified items\n")
				outputContent, err = preservation.MergeWithPreservation(basePath, sourcePath, tablesConfig, tpl)
				if err != nil {
					log.Fatalf("❌ Error merging with preservation for job '%s': %v", jobName, err)
				}
			} else {
				// Mode: Only specified tables (current behavior)
				results, err := merger.MergeTables(basePath, sourcePath, tablesConfig)
				if err != nil {
					log.Fatalf("❌ Error merging job '%s': %v", jobName, err)
				}

				var buf []byte
				for j, result := range results {
					if j > 0 {
						buf = append(buf, []byte("\n\n")...)
					}

					var resultBuf bytes.Buffer
					if err := tpl.Execute(&resultBuf, result); err != nil {
						log.Fatalf("❌ Error executing template for table '%s': %v", result.TableName, err)
					}
					buf = append(buf, resultBuf.Bytes()...)
				}
				outputContent = string(buf)
			}

			// Write output file
			if err := os.WriteFile(outputPath, []byte(outputContent), 0644); err != nil {
				log.Fatalf("❌ Error writing output file '%s': %v", outputPath, err)
			}

			fmt.Printf("  ✓ Base: %s\n", filepath.Base(basePath))
			fmt.Printf("  ✓ Source: %s\n", filepath.Base(sourcePath))
			fmt.Printf("  ✓ Output: %s\n", outputPath)
			fmt.Printf("  ✓ Tables: %d\n\n", len(job.Tables))
		}

		fmt.Printf("🎉 All %d job(s) processed successfully!\n", len(settings.Jobs))
	},
}

func init() {
	rootCmd.Flags().StringVarP(&inputDir, "inputs", "i", "input", "Input directory containing settings.json")
	rootCmd.SetVersionTemplate(fmt.Sprintf("v%s\n", version))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
