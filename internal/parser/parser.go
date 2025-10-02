package parser

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/yuin/gopher-lua/ast"
	"github.com/yuin/gopher-lua/parse"
)

// Parse parses a Lua file and extracts a specific table by name.
// Returns the parsed Table structure or an error if parsing fails.
func Parse(
	reader io.Reader,
	sourceName string,
	tableName string,
) (*Table, error) {
	chunk, err := parse.Parse(reader, sourceName)
	if err != nil {
		return nil, fmt.Errorf("parser.Parse: %w", err)
	}

	tableNode, err := findTable(chunk, tableName)
	if err != nil {
		return nil, err
	}

	return parseTable(tableNode)
}

// findTable searches for a table assignment in the AST by name
func findTable(chunk []ast.Stmt, tableName string) (*ast.TableExpr, error) {
	for _, stmt := range chunk {
		if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
			ident, ok := assignStmt.Lhs[0].(*ast.IdentExpr)
			if ok && ident.Value == tableName {
				if tableNode, ok := assignStmt.Rhs[0].(*ast.TableExpr); ok {
					return tableNode, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("parser.findTable: table '%s' was not found", tableName)
}

// parseTable converts an AST table expression into our Table structure
func parseTable(node *ast.TableExpr) (*Table, error) {
	table := NewTable()
	for _, field := range node.Fields {
		value, err := parseValue(field.Value)
		if err != nil {
			return nil, err
		}

		if field.Key == nil {
			table.AddOrReplace("", value)
		} else {
			key, err := parseKey(field.Key)
			if err != nil {
				return nil, fmt.Errorf("failed to convert AST key: %w", err)
			}
			table.AddOrReplace(key, value)
		}
	}
	return table, nil
}

// valueOf extracts the string value from an AST expression
func valueOf(exp ast.Expr) (string, error) {
	switch v := exp.(type) {
	case *ast.IdentExpr:
		return v.Value, nil
	case *ast.StringExpr:
		return v.Value, nil
	case *ast.NumberExpr:
		return v.Value, nil
	default:
		return "", fmt.Errorf("parser.valueOf: unsupported AST key type: %T", exp)
	}
}

// parseValue converts an AST value expression into our Value structure
func parseValue(exp ast.Expr) (*Value, error) {
	switch v := exp.(type) {
	case *ast.NilExpr:
		return &Value{Type: TypeNil, value: nil}, nil
	case *ast.TrueExpr:
		return &Value{Type: TypeBoolean, value: true}, nil
	case *ast.FalseExpr:
		return &Value{Type: TypeBoolean, value: false}, nil
	case *ast.StringExpr:
		return &Value{Type: TypeString, value: v.Value}, nil
	case *ast.NumberExpr:
		num, err := strconv.ParseFloat(v.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert AST number '%s': %w", v.Value, err)
		}
		return &Value{Type: TypeNumber, value: num}, nil
	case *ast.TableExpr:
		table, err := parseTable(v)
		if err != nil {
			return nil, err
		}
		return &Value{Type: TypeTable, value: table}, nil
	case *ast.FunctionExpr:
		return &Value{Type: TypeFunction, value: "<function>"}, nil
	case *ast.AttrGetExpr:
		obj, err := valueOf(v.Object)
		if err != nil {
			return nil, err
		}
		key, err := valueOf(v.Key)
		if err != nil {
			return nil, err
		}

		str := fmt.Sprintf("%s.%s", obj, key)
		return &Value{Type: TypeVariable, value: str}, nil
	case *ast.IdentExpr:
		return &Value{Type: TypeVariable, value: v.Value}, nil
	default:
		return nil, fmt.Errorf("unsupported AST value type: %T", exp)
	}
}

// parseKey converts an AST key expression into a string representation
func parseKey(exp ast.Expr) (string, error) {
	switch v := exp.(type) {
	case *ast.IdentExpr:
		return v.Value, nil
	case *ast.StringExpr:
		if strings.Contains(v.Value, ".") {
			return fmt.Sprintf("[\"%s\"]", v.Value), nil
		}
		return v.Value, nil
	case *ast.NumberExpr:
		return fmt.Sprintf("[%s]", v.Value), nil
	case *ast.AttrGetExpr:
		obj, err := valueOf(v.Object)
		if err != nil {
			return "", err
		}
		key, err := valueOf(v.Key)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("[%s.%s]", obj, key), nil
	default:
		return "", fmt.Errorf("unsupported AST key type: %T", exp)
	}
}
