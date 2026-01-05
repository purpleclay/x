package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnumInt(t *testing.T) {
	type TrustLevel int
	const (
		TrustUnknown TrustLevel = iota + 1
		TrustNever
		TrustMarginal
	)

	e := Enum(TrustUnknown, TrustUnknown, TrustNever, TrustMarginal)

	assert.Equal(t, "int", e.BaseType())
	assert.Equal(t, TrustUnknown, e.Get())
	assert.Equal(t, "1", e.String())
	assert.Equal(t, "1|2|3", e.Type())
}

func TestEnumString(t *testing.T) {
	type Format string
	const (
		FormatJSON Format = "json"
		FormatYAML Format = "yaml"
	)

	e := Enum(FormatJSON, FormatJSON, FormatYAML)

	assert.Equal(t, "string", e.BaseType())
	assert.Equal(t, FormatJSON, e.Get())
	assert.Equal(t, "json", e.String())
	assert.Equal(t, "json|yaml", e.Type())
}

func TestEnumStringWithHelp(t *testing.T) {
	type Format string
	const (
		FormatJSON Format = "json"
		FormatYAML Format = "yaml"
	)

	e := Enum(FormatJSON, FormatJSON, FormatYAML).
		WithHelp("JavaScript Object Notation", "YAML Ain't Markup Language")

	assert.True(t, e.HasHelp())

	entries := e.HelpEntries()
	require.Len(t, entries, 2)
	assert.Equal(t, "json", entries[0].Name)
	assert.Equal(t, "JavaScript Object Notation", entries[0].Help)
	assert.Equal(t, "yaml", entries[1].Name)
	assert.Equal(t, "YAML Ain't Markup Language", entries[1].Help)
}

func TestEnumSetFailsWithUnmatchedValue(t *testing.T) {
	type Format string
	const (
		FormatJSON Format = "json"
		FormatYAML Format = "yaml"
	)

	e := Enum(FormatJSON, FormatJSON, FormatYAML)

	err := e.Set("xml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be one of")
	assert.Equal(t, FormatJSON, e.Get()) // unchanged
}
