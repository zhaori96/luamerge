package merger

import (
	"fmt"
	"luamerge/internal/parser"
	"os"
)

// applyRules recursively applies merge rules to a table.
// Supports deep merging at any nesting level.
func applyRules(base, source *parser.Table, rules map[string]any) {
	for ruleKey, ruleValue := range rules {
		baseValue, baseExists := base.Get(ruleKey)
		sourceValue, sourceExists := source.Get(ruleKey)

		// If the key doesn't exist in both, skip
		if !baseExists || !sourceExists {
			continue
		}

		// If the rule is true, replace the value completely
		if ruleBool, ok := ruleValue.(bool); ok && ruleBool {
			base.AddOrReplace(ruleKey, sourceValue)
			continue
		}

		// If the rule is a map (object), do recursive merge
		if nestedRules, ok := ruleValue.(map[string]any); ok {
			baseTable, baseIsTable := baseValue.Table()
			sourceTable, sourceIsTable := sourceValue.Table()

			// Both need to be tables for recursive merge
			if baseIsTable == nil && sourceIsTable == nil {
				applyRules(baseTable, sourceTable, nestedRules)
			}
		}
	}
}

// mergeInternal applies rules to all entries of a top-level table.
// Iterates over all entries in the base table and applies merge rules.
func mergeInternal(baseTable, sourceTable *parser.Table, rules map[string]any) {
	// If rules is nil or empty, replace the entire table
	if len(rules) == 0 {
		// Clear the base table and copy everything from source
		for sourceEntry := range sourceTable.Range() {
			baseTable.AddOrReplace(sourceEntry.Name, sourceEntry.Value)
		}
		return
	}

	// If the rules indicate complete replacement (true), replace everything
	for key, ruleValue := range rules {
		if ruleBool, ok := ruleValue.(bool); ok && ruleBool {
			// Complete replacement of all entries
			for sourceEntry := range sourceTable.Range() {
				baseTable.AddOrReplace(sourceEntry.Name, sourceEntry.Value)
			}
			return
		}

		// If it's a map of rules, apply recursively to each entry
		if nestedRules, ok := ruleValue.(map[string]any); ok {
			for baseEntry := range baseTable.Range() {
				sourceEntry, sourceExists := sourceTable.Get(baseEntry.Name)
				if !sourceExists {
					continue
				}

				baseSubTable, baseIsTable := baseEntry.Value.Table()
				sourceSubTable, sourceIsTable := sourceEntry.Table()

				if baseIsTable == nil && sourceIsTable == nil {
					applyRules(baseSubTable, sourceSubTable, nestedRules)
				}
			}
		}

		_ = key
	}
}

// MergeTables merges multiple tables from two Lua files.
// Receives the file paths and a table configuration map.
// Returns a slice of Result containing the merged tables.
func MergeTables(basePath, sourcePath string, tablesConfig map[string]map[string]any) ([]Result, error) {
	// Input validations
	if basePath == "" {
		return nil, fmt.Errorf("base file path cannot be empty")
	}
	if sourcePath == "" {
		return nil, fmt.Errorf("source file path cannot be empty")
	}
	if len(tablesConfig) == 0 {
		return nil, fmt.Errorf("tables configuration cannot be empty")
	}

	// Check if files exist
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("base file not found: %s", basePath)
	}
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("source file not found: %s", sourcePath)
	}

	baseF, err := os.Open(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open base file '%s': %w", basePath, err)
	}
	defer baseF.Close()

	sourceF, err := os.Open(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file '%s': %w", sourcePath, err)
	}
	defer sourceF.Close()

	var results []Result

	for tableName, fieldsToReplace := range tablesConfig {
		if tableName == "" {
			return nil, fmt.Errorf("empty table name found in configuration")
		}

		baseTable, err := parser.Parse(baseF, basePath, tableName)
		if err != nil {
			return nil, fmt.Errorf("failed to parse table '%s' in base file: %w", tableName, err)
		}
		baseF.Seek(0, 0)

		sourceTable, err := parser.Parse(sourceF, sourcePath, tableName)
		if err != nil {
			return nil, fmt.Errorf("failed to parse table '%s' in source file: %w", tableName, err)
		}
		sourceF.Seek(0, 0)

		mergeInternal(baseTable, sourceTable, fieldsToReplace)

		results = append(results, Result{
			TableName: tableName,
			Table:     baseTable,
		})
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no tables were processed")
	}

	return results, nil
}
