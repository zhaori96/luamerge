package parser

import (
	"errors"
	"fmt"
	"iter"
	"strconv"
	"strings"
)

// Type represents the type of a Lua value
type Type int

const (
	TypeNil Type = iota
	TypeBoolean
	TypeNumber
	TypeString
	TypeTable
	TypeFunction
	TypeVariable
)

// Table represents a Lua table with named or indexed values
type Table struct {
	values       []*NamedValue
	index        map[string]int
	currentIndex int
}

// NewTable creates a new empty Table
func NewTable() *Table {
	return &Table{
		index: make(map[string]int),
	}
}

// AddOrReplace adds a new value to the table or replaces an existing one
func (t *Table) AddOrReplace(name string, value *Value) {
	if current, ok := t.Get(name); ok {
		current.value = value
		return
	}

	if name == "" {
		t.currentIndex++
		name = fmt.Sprintf("[%d]", t.currentIndex)
	}

	t.values = append(t.values, &NamedValue{name, value})
	t.index[name] = len(t.values) - 1
}

// Get retrieves a value from the table by key
func (t *Table) Get(key string) (*Value, bool) {
	// First try to find by exact key match
	index, ok := t.index[key]
	if ok && index < len(t.values) {
		return t.values[index].Value, true
	}

	// If not found and the key looks numeric, try to access by index
	// Remove brackets if present (e.g., "[123]" becomes "123")
	trimmedKey := strings.TrimPrefix(strings.TrimSuffix(key, "]"), "[")
	numIndex, err := strconv.Atoi(trimmedKey)
	if err != nil {
		return nil, false
	}

	// Search for the formatted index key
	formattedKey := fmt.Sprintf("[%d]", numIndex)
	if idx, ok := t.index[formattedKey]; ok && idx < len(t.values) {
		return t.values[idx].Value, true
	}

	return nil, false
}

// Range returns an iterator over all named values in the table
func (t *Table) Range() iter.Seq[*NamedValue] {
	return func(yield func(*NamedValue) bool) {
		for _, value := range t.values {
			if !yield(value) {
				return
			}
		}
	}
}

// NamedValue represents a table entry with a name and value
type NamedValue struct {
	Name  string
	Value *Value
}

// Value represents a Lua value with its type
type Value struct {
	Type  Type
	value any
}

// Value returns the underlying Go value
func (v *Value) Value() any {
	return v.value
}

// String returns the string value if the type is TypeString
func (v *Value) String() (string, error) {
	if v.Type != TypeString {
		return "", errors.New("is not a string")
	}

	if value, ok := v.value.(*Value); ok {
		return value.String()
	}

	value, ok := v.value.(string)
	if !ok {
		return "", fmt.Errorf("value is of type %T, not string", v.value)
	}

	value = strings.ReplaceAll(value, "\n", "\\n")
	value = strings.ReplaceAll(value, "\t", "\\t")
	value = strings.Trim(value, "\"")
	value = strings.ReplaceAll(value, `\"`, "\"")
	value = strings.ReplaceAll(value, `"`, `\"`)
	return value, nil
}

// Number returns the numeric value if the type is TypeNumber
func (v *Value) Number() (float64, error) {
	if v.Type != TypeNumber {
		return 0, errors.New("is not a number")
	}

	if value, ok := v.value.(*Value); ok {
		return value.Number()
	}

	value, ok := v.value.(float64)
	if !ok {
		return 0, fmt.Errorf("value is of type %T, not float64", v.value)
	}

	return value, nil
}

// Boolean returns the boolean value if the type is TypeBoolean
func (v *Value) Boolean() (bool, error) {
	if v.Type != TypeBoolean {
		return false, errors.New("is not a boolean")
	}

	if value, ok := v.value.(*Value); ok {
		return value.Boolean()
	}

	value, ok := v.value.(bool)
	if !ok {
		return false, fmt.Errorf("value is of type %T, not bool", v.value)
	}
	return value, nil
}

// Table returns the table value if the type is TypeTable
func (v *Value) Table() (*Table, error) {
	if v.Type != TypeTable {
		return nil, errors.New("is not a table")
	}

	if value, ok := v.value.(*Value); ok {
		return value.Table()
	}

	value, ok := v.value.(*Table)
	if !ok {
		return nil, fmt.Errorf("value is of type %T, not *Table", v.value)
	}

	return value, nil
}

// Variable returns the variable name if the type is TypeVariable
func (v *Value) Variable() (string, error) {
	if v.Type != TypeVariable {
		return "", errors.New("is not a variable")
	}

	if value, ok := v.value.(*Value); ok {
		return value.Variable()
	}

	value, ok := v.value.(string)
	if !ok {
		return "", fmt.Errorf("value is of type %T, not string", v.value)
	}
	return value, nil
}
