package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractTextBetweenToPatterns(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		fromPattern string
		toPattern   string
		wantValue   string
	}{
		{
			name:        "simple extraction",
			source:      "hello world from here to there",
			fromPattern: "from",
			toPattern:   "to",
			wantValue:   " here ",
		},
		{
			name: "multi lines extraction",
			source: `hello world
	from here to there
`,
			fromPattern: "world",
			toPattern:   "from",
			wantValue: `
	`,
		},
		{
			name:        "no from pattern",
			source:      "hello world to there",
			fromPattern: "xxx",
			toPattern:   "to",
			wantValue:   "hello world ",
		},
		{
			name:        "no to pattern",
			source:      "hello world from here",
			fromPattern: "from",
			toPattern:   "xxx",
			wantValue:   "",
		},
		{
			name:        "from pattern at start",
			source:      "from here to there",
			fromPattern: "from",
			toPattern:   "to",
			wantValue:   " here ",
		},
		{
			name:        "to pattern at end",
			source:      "hello world from here to",
			fromPattern: "from",
			toPattern:   "to",
			wantValue:   " here ",
		},
		{
			name:        "multiple from patterns",
			source:      "hello world from here from there to",
			fromPattern: "from",
			toPattern:   "to",
			wantValue:   " there ",
		},
		{
			name:        "multiple to patterns",
			source:      "hello world from here to there to",
			fromPattern: "from",
			toPattern:   "to",
			wantValue:   " here ",
		},
		{
			name:        "empty source",
			source:      "",
			fromPattern: "from",
			toPattern:   "to",
			wantValue:   "",
		},
		{
			name:        "empty from pattern",
			source:      "hello world to there",
			fromPattern: "",
			toPattern:   "to",
			wantValue:   "",
		},
		{
			name:        "empty to pattern",
			source:      "hello world from here",
			fromPattern: "from",
			toPattern:   "",
			wantValue:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantValue, ExtractTextBetweenToPatterns(tt.source, tt.fromPattern, tt.toPattern))
		})
	}
}

func TestReplaceLast(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		pattern     string
		replacement string
		wantValue   string
	}{
		{
			name:        "simple replacement",
			source:      "hello world, world!",
			pattern:     "world",
			replacement: "earth",
			wantValue:   "hello world, earth!",
		},
		{
			name:        "no pattern",
			source:      "hello world, world!",
			pattern:     "xxx",
			replacement: "earth",
			wantValue:   "hello world, world!",
		},
		{
			name:        "empty replacement",
			source:      "hello world, world!",
			pattern:     "world",
			replacement: "",
			wantValue:   "hello world, !",
		},
		{
			name:        "pattern at start",
			source:      "world, world!",
			pattern:     "world",
			replacement: "earth",
			wantValue:   "world, earth!",
		},
		{
			name:        "pattern at end",
			source:      "hello world, world",
			pattern:     "world",
			replacement: "earth",
			wantValue:   "hello world, earth",
		},
		{
			name:        "multiple patterns",
			source:      "hello world, world, world!",
			pattern:     "world",
			replacement: "earth",
			wantValue:   "hello world, world, earth!",
		},
		{
			name:        "empty source",
			source:      "",
			pattern:     "world",
			replacement: "earth",
			wantValue:   "",
		},
		{
			name:        "empty pattern",
			source:      "hello world, world!",
			pattern:     "",
			replacement: "earth",
			wantValue:   "hello world, world!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantValue, ReplaceLast(tt.source, tt.pattern, tt.replacement))
		})
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		source    string
		wantValue string
	}{
		{
			source:    "Tenant",
			wantValue: "Tenants",
		},
		{
			source:    "AdvancedPayment",
			wantValue: "AdvancedPayments",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.wantValue, Pluralize(tt.source))
		})
	}
}

func TestToUpperCaseFirstChar(t *testing.T) {
	tests := []struct {
		source    string
		wantValue string
	}{
		{
			source:    "hello",
			wantValue: "Hello",
		},
		{
			source:    "pcbCalculator",
			wantValue: "PcbCalculator",
		},
		{
			source:    "guid",
			wantValue: "Guid",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.wantValue, ToUpperCaseFirstChar(tt.source))
		})
	}
}

func TestSplitCamelCase(t *testing.T) {
	tests := []struct {
		source    string
		wantValue string
	}{
		{
			source:    "AdvancedPayment",
			wantValue: "Advanced Payment",
		},
		{
			source:    "Tenant",
			wantValue: "Tenant",
		},
		{
			source:    "PMRBTable",
			wantValue: "PMRB Table",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.wantValue, SplitCamelCase(tt.source, " "))
		})
	}
}
