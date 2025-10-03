package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GlobalOptions represents global options for all jobs
type GlobalOptions struct {
	KeepUnmergedItems bool `json:"keepUnmergedItems"`
}

// JobOptions represents job-specific options (can override global options)
type JobOptions struct {
	KeepUnmergedItems *bool `json:"keepUnmergedItems,omitempty"`
}

// Job represents a merge task configured in settings.json
type Job struct {
	Name    string                 `json:"name"`
	Base    string                 `json:"base"`
	Source  string                 `json:"source"`
	Output  string                 `json:"output"`
	Tables  map[string]any         `json:"tables"`
	Options *JobOptions            `json:"options,omitempty"`
}

// GetTablesConfig normalizes the tables configuration to the format expected by the merger
func (j *Job) GetTablesConfig() map[string]map[string]any {
	result := make(map[string]map[string]any)

	for tableName, value := range j.Tables {
		switch v := value.(type) {
		case bool:
			// If true, replace everything (represented by empty map or nil)
			if v {
				result[tableName] = nil
			}
		case map[string]any:
			// Already a map, use directly
			result[tableName] = v
		default:
			// Ignore invalid values
			continue
		}
	}

	return result
}

// Settings represents the complete job configuration
type Settings struct {
	Options *GlobalOptions `json:"options,omitempty"`
	Jobs    []Job          `json:"jobs"`
}

// GetKeepUnmergedItems returns whether to keep unmerged items, respecting the hierarchy
func (j *Job) GetKeepUnmergedItems(globalOptions *GlobalOptions) bool {
	// Job options take priority over global
	if j.Options != nil && j.Options.KeepUnmergedItems != nil {
		return *j.Options.KeepUnmergedItems
	}

	// Use global if configured
	if globalOptions != nil {
		return globalOptions.KeepUnmergedItems
	}

	// Default: false (current behavior)
	return false
}

// LoadSettingsFromInput loads the settings.json file from the input folder
func LoadSettingsFromInput(inputDir string) (*Settings, error) {
	settingsPath := filepath.Join(inputDir, "settings.json")

	// Check if the file exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("settings.json file not found in: %s", inputDir)
	}

	b, err := os.ReadFile(settingsPath)
	if err != nil {
		return nil, fmt.Errorf("error reading settings.json: %w", err)
	}

	if len(b) == 0 {
		return nil, fmt.Errorf("settings.json file is empty: %s", settingsPath)
	}

	var settings Settings
	if err := json.Unmarshal(b, &settings); err != nil {
		return nil, fmt.Errorf("error parsing settings.json: %w", err)
	}

	if len(settings.Jobs) == 0 {
		return nil, fmt.Errorf("no jobs configured in settings.json")
	}

	// Validate each job
	for i, job := range settings.Jobs {
		if err := validateJob(job, i); err != nil {
			return nil, err
		}
	}

	return &settings, nil
}

// validateJob validates an individual job
func validateJob(job Job, index int) error {
	jobID := fmt.Sprintf("job[%d]", index)
	if job.Name != "" {
		jobID = fmt.Sprintf("job[%d] (%s)", index, job.Name)
	}

	if job.Base == "" {
		return fmt.Errorf("%s: 'base' field is required", jobID)
	}

	if job.Source == "" {
		return fmt.Errorf("%s: 'source' field is required", jobID)
	}

	if job.Output == "" {
		return fmt.Errorf("%s: 'output' field is required", jobID)
	}

	if len(job.Tables) == 0 {
		return fmt.Errorf("%s: 'tables' field is required and must contain at least one table", jobID)
	}

	return nil
}

// ResolveJobPaths resolves the relative paths of a job based on the input folder
func ResolveJobPaths(job Job, inputDir string) (basePath, sourcePath, outputPath string, err error) {
	// Resolve base and source relative to the input folder
	basePath = filepath.Join(inputDir, job.Base)
	sourcePath = filepath.Join(inputDir, job.Source)

	// Resolve output
	if filepath.IsAbs(job.Output) {
		// Absolute path
		outputPath = job.Output
	} else if filepath.Dir(job.Output) == "." {
		// Just the filename, goes to output/ next to input/
		outputDir := filepath.Join(filepath.Dir(inputDir), "output")
		outputPath = filepath.Join(outputDir, job.Output)
	} else {
		// Relative path (may include ../)
		outputPath = filepath.Join(inputDir, job.Output)
		outputPath = filepath.Clean(outputPath)
	}

	// Validate that input files exist
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return "", "", "", fmt.Errorf("base file not found: %s", basePath)
	}

	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return "", "", "", fmt.Errorf("source file not found: %s", sourcePath)
	}

	return basePath, sourcePath, outputPath, nil
}
