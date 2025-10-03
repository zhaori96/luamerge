package preservation

import (
	"bytes"
	"fmt"
	"luamerge/internal/merger"
	"os"
	"regexp"
	"text/template"
)

// TablePosition represents the position of a table in the file
type TablePosition struct {
	Name      string
	StartLine int
	EndLine   int
	StartPos  int
	EndPos    int
}

// FindTablePositions finds the positions of all tables in the text
func FindTablePositions(content string, tableNames []string) (map[string]*TablePosition, error) {
	positions := make(map[string]*TablePosition)

	for _, tableName := range tableNames {
		// Pattern to find table declaration: TableName = {
		pattern := fmt.Sprintf(`(?m)^%s\s*=\s*\{`, regexp.QuoteMeta(tableName))
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex for table '%s': %w", tableName, err)
		}

		// Find the starting position of the table
		matches := re.FindStringIndex(content)
		if matches == nil {
			continue // Table not found
		}

		startPos := matches[0]

		// Find the end of the table by balancing braces
		endPos := findTableEndInContent(content, matches[1]-1) // -1 to start from {

		if endPos == -1 {
			continue // Could not find end
		}

		positions[tableName] = &TablePosition{
			Name:      tableName,
			StartLine: 0, // No longer used
			EndLine:   0, // No longer used
			StartPos:  startPos,
			EndPos:    endPos,
		}
	}

	return positions, nil
}

// findTableEndInContent finds the end of a table in the content by balancing braces
func findTableEndInContent(content string, startPos int) int {
	depth := 0
	inString := false
	escapeNext := false

	for i := startPos; i < len(content); i++ {
		char := content[i]

		// Handle escape sequences in strings
		if escapeNext {
			escapeNext = false
			continue
		}

		if char == '\\' {
			escapeNext = true
			continue
		}

		// Handle strings (ignore braces inside strings)
		if char == '"' {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		// Balance braces
		switch char {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				// Found the end - include the line up to newline
				for j := i + 1; j < len(content); j++ {
					if content[j] == '\n' {
						return j + 1
					}
				}
				return i + 1
			}
		}
	}

	return -1 // Did not find end
}

// ReplaceTablesInText replaces the merged tables in the original text
func ReplaceTablesInText(originalContent string, mergedResults []merger.Result, tpl *template.Template) (string, error) {
	// Extract table names
	tableNames := make([]string, len(mergedResults))
	for i, result := range mergedResults {
		tableNames[i] = result.TableName
	}

	// Find table positions
	positions, err := FindTablePositions(originalContent, tableNames)
	if err != nil {
		return "", err
	}

	// Sort positions by StartPos (from end to beginning, to avoid affecting offsets)
	sortedPositions := make([]*TablePosition, 0, len(positions))
	for _, pos := range positions {
		sortedPositions = append(sortedPositions, pos)
	}

	// Bubble sort in reverse (from end to beginning)
	for i := 0; i < len(sortedPositions)-1; i++ {
		for j := 0; j < len(sortedPositions)-i-1; j++ {
			if sortedPositions[j].StartPos < sortedPositions[j+1].StartPos {
				sortedPositions[j], sortedPositions[j+1] = sortedPositions[j+1], sortedPositions[j]
			}
		}
	}

	// Create map of results by name
	resultMap := make(map[string]merger.Result)
	for _, result := range mergedResults {
		resultMap[result.TableName] = result
	}

	// Replace from end to beginning (to avoid affecting offsets)
	result := originalContent
	for _, pos := range sortedPositions {
		mergedResult, exists := resultMap[pos.Name]
		if !exists {
			continue
		}

		// Generate Lua from merged table using template
		var buf bytes.Buffer
		if err := tpl.Execute(&buf, mergedResult); err != nil {
			return "", fmt.Errorf("error generating Lua for table '%s': %w", pos.Name, err)
		}

		// Replace the section
		before := result[:pos.StartPos]
		after := result[pos.EndPos:]
		result = before + buf.String() + after
	}

	return result, nil
}

// MergeWithPreservation performs merge while preserving unspecified items
func MergeWithPreservation(basePath, sourcePath string, tablesConfig map[string]map[string]any, tpl *template.Template) (string, error) {
	// Read base file as text
	baseContent, err := os.ReadFile(basePath)
	if err != nil {
		return "", fmt.Errorf("error reading base file: %w", err)
	}

	// Perform normal merge of specified tables
	mergedResults, err := merger.MergeTables(basePath, sourcePath, tablesConfig)
	if err != nil {
		return "", err
	}

	// Replace tables in original text
	result, err := ReplaceTablesInText(string(baseContent), mergedResults, tpl)
	if err != nil {
		return "", err
	}

	return result, nil
}
