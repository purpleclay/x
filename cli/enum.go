package cli

import (
	"fmt"
	"reflect"
	"strings"
)

// Enumerable defines the constraint for enum types.
type Enumerable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~string
}

// EnumOption represents an enum value with an optional help description.
type EnumOption struct {
	Name string
	Help string
}

// EnumHelper is implemented by enum values that have help text for their options.
type EnumHelper interface {
	HasHelp() bool
	HelpEntries() []EnumOption
	BaseType() string
}

// EnumValue implements pflag.Value for type-safe enumeration flags.
type EnumValue[T Enumerable] struct {
	value    T
	names    map[T]string
	values   map[string]T
	allowed  []string
	help     map[string]string
	baseType string
}

// Enum creates a new type-safe enum flag. The first argument is the default
// value, followed by all allowed values. For string-based enums, the string
// value is used as the display name. For integer-based enums, the numeric
// value is used as the display name.
//
// String enum example:
//
//	type Format string
//
//	const (
//	    FormatJSON Format = "json"
//	    FormatYAML Format = "yaml"
//	    FormatTOML Format = "toml"
//	)
//
//	format := cli.Enum(FormatJSON, FormatJSON, FormatYAML, FormatTOML)
//	cmd.Flags().VarP(format, "format", "f", "output format")
//
// Integer enum example:
//
//	type TrustLevel int
//
//	const (
//	    TrustUnknown  TrustLevel = iota + 1  // 1
//	    TrustNever                           // 2
//	    TrustMarginal                        // 3
//	)
//
//	trust := cli.Enum(TrustUnknown, TrustUnknown, TrustNever, TrustMarginal)
//	cmd.Flags().Var(trust, "trust-level", "GPG trust level")
func Enum[T Enumerable](def T, allowed ...T) *EnumValue[T] {
	names := make(map[T]string, len(allowed))
	values := make(map[string]T, len(allowed))
	orderedNames := make([]string, 0, len(allowed))

	for _, v := range allowed {
		name := fmt.Sprintf("%v", v)
		names[v] = name
		values[name] = v
		orderedNames = append(orderedNames, name)
	}

	baseType := "int"
	if reflect.TypeOf(def).Kind() == reflect.String {
		baseType = "string"
	}

	return &EnumValue[T]{
		value:    def,
		names:    names,
		values:   values,
		allowed:  orderedNames,
		baseType: baseType,
	}
}

// WithHelp adds help text for each enum value in order. The help strings
// correspond to the enum values in the order they were defined.
//
//	format := cli.Enum(FormatJSON, FormatJSON, FormatYAML, FormatTOML).
//	    WithHelp(
//	        "JavaScript Object Notation",
//	        "YAML Ain't Markup Language",
//	        "Tom's Obvious Minimal Language",
//	    )
func (e *EnumValue[T]) WithHelp(help ...string) *EnumValue[T] {
	if len(help) == 0 {
		return e
	}

	e.help = make(map[string]string, len(help))
	for i, h := range help {
		if i < len(e.allowed) && h != "" {
			e.help[e.allowed[i]] = h
		}
	}

	return e
}

// String returns the string representation of the current value.
func (e *EnumValue[T]) String() string {
	if name, ok := e.names[e.value]; ok {
		return name
	}
	return ""
}

// Set validates and sets the value from a string.
func (e *EnumValue[T]) Set(s string) error {
	if v, ok := e.values[s]; ok {
		e.value = v
		return nil
	}
	return fmt.Errorf("must be one of: %s", strings.Join(e.allowed, ", "))
}

// Type returns the type name for help output, showing all allowed values.
func (e *EnumValue[T]) Type() string {
	return strings.Join(e.allowed, "|")
}

// Get returns the current typed enum value.
//
//nolint:ireturn
func (e *EnumValue[T]) Get() T {
	return e.value
}

// HasHelp returns true if this enum has help text for its values.
func (e *EnumValue[T]) HasHelp() bool {
	return len(e.help) > 0
}

// HelpEntries returns the enum values with their help text in display order.
func (e *EnumValue[T]) HelpEntries() []EnumOption {
	entries := make([]EnumOption, len(e.allowed))
	for i, name := range e.allowed {
		entries[i] = EnumOption{
			Name: name,
			Help: e.help[name],
		}
	}
	return entries
}

// BaseType returns the underlying type name ("string" or "int").
func (e *EnumValue[T]) BaseType() string {
	return e.baseType
}
