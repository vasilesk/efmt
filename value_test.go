package efmt_test

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"testing"

	"github.com/vasilesk/efmt"
)

func TestValuePairs(t *testing.T) {
	// Create a chain of errors with key-values
	baseErr := efmt.New("base error", efmt.KV("base", "value"), efmt.KV("common", "base"))
	middleErr := efmt.Wrap(baseErr, "middle error", efmt.KV("middle", 42), efmt.KV("common", "middle"))
	topErr := efmt.Wrap(middleErr, "top error", efmt.KV("top", true), efmt.KV("common", "top"))

	// Test valuePairs for different levels of the error chain
	tests := []struct {
		name      string
		err       error
		wantCount int
		wantKeys  []string
	}{
		{
			name:      "nil error",
			err:       nil,
			wantCount: 0,
			wantKeys:  []string{},
		},
		{
			name:      "base error",
			err:       baseErr,
			wantCount: 2,
			wantKeys:  []string{"base", "common"},
		},
		{
			name:      "middle error",
			err:       middleErr,
			wantCount: 4,
			wantKeys:  []string{"middle", "common", "base", "common"},
		},
		{
			name:      "top error",
			err:       topErr,
			wantCount: 6,
			wantKeys:  []string{"top", "common", "middle", "common", "base", "common"},
		},
		{
			name:      "standard error",
			err:       errors.New("standard error"),
			wantCount: 0,
			wantKeys:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pairs := efmt.ValuePairs(tt.err)

			// Check count
			if len(pairs) != tt.wantCount {
				t.Errorf("ValuePairs() count = %v, want %v", len(pairs), tt.wantCount)

				return
			}

			// Check keys (order matters)
			if len(pairs) > 0 {
				for i, pair := range pairs {
					if i >= len(tt.wantKeys) {
						t.Errorf("ValuePairs() unexpected key at index %d: %s", i, pair.Key)

						continue
					}

					if pair.Key != tt.wantKeys[i] {
						t.Errorf("ValuePairs() key at index %d = %s, want %s", i, pair.Key, tt.wantKeys[i])
					}
				}
			}
		})
	}
}

func TestValues(t *testing.T) {
	// Create a chain of errors with key-values
	baseErr := efmt.New("base error", efmt.KV("base", "value"), efmt.KV("common", "base"))
	middleErr := efmt.Wrap(baseErr, "middle error", efmt.KV("middle", 42), efmt.KV("common", "middle"))
	topErr := efmt.Wrap(middleErr, "top error", efmt.KV("top", true), efmt.KV("common", "top"))

	tests := []struct {
		name      string
		err       error
		wantCount int
		wantMap   map[string]any
	}{
		{
			name:      "nil error",
			err:       nil,
			wantCount: 0,
			wantMap:   map[string]any{},
		},
		{
			name:      "base error",
			err:       baseErr,
			wantCount: 2,
			wantMap: map[string]any{
				"base":   "value",
				"common": "base",
			},
		},
		{
			name:      "middle error",
			err:       middleErr,
			wantCount: 3,
			wantMap: map[string]any{
				"base":   "value",
				"middle": 42,
				"common": "middle", // Should override base's "common" value
			},
		},
		{
			name:      "top error",
			err:       topErr,
			wantCount: 4,
			wantMap: map[string]any{
				"base":   "value",
				"middle": 42,
				"top":    true,
				"common": "top", // Should override middle and base's "common" value
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := efmt.Values(tt.err)

			// Check count
			if len(values) != tt.wantCount {
				t.Errorf("Values() count = %v, want %v", len(values), tt.wantCount)
			}

			// Check key-value pairs
			for k, wantV := range tt.wantMap {
				if gotV, ok := values[k]; !ok {
					t.Errorf("Values() missing key %q", k)
				} else if gotV != wantV {
					t.Errorf("Values()[%q] = %v, want %v", k, gotV, wantV)
				}
			}

			// Check for unexpected keys
			for k := range values {
				if _, ok := tt.wantMap[k]; !ok {
					t.Errorf("Values() unexpected key %q", k)
				}
			}
		})
	}
}

func TestValue(t *testing.T) {
	// Create an error with various value types
	err := efmt.New("test error",
		efmt.KV("string", "string value"),
		efmt.KV("int", 42),
		efmt.KV("bool", true),
		efmt.KV("float", 3.14),
	)

	// Create a chain to test value overrides
	wrappedErr := efmt.Wrap(err, "wrapped", efmt.KV("string", "overridden"))

	tests := []struct {
		name      string
		err       error
		key       string
		wantValue any
		wantOk    bool
	}{
		{
			name:      "get string value",
			err:       err,
			key:       "string",
			wantValue: "string value",
			wantOk:    true,
		},
		{
			name:      "get int value",
			err:       err,
			key:       "int",
			wantValue: 42,
			wantOk:    true,
		},
		{
			name:      "get bool value",
			err:       err,
			key:       "bool",
			wantValue: true,
			wantOk:    true,
		},
		{
			name:      "get float value",
			err:       err,
			key:       "float",
			wantValue: 3.14,
			wantOk:    true,
		},
		{
			name:      "overridden value in chain",
			err:       wrappedErr,
			key:       "string",
			wantValue: "overridden",
			wantOk:    true,
		},
		{
			name:      "non-existent key",
			err:       err,
			key:       "nonexistent",
			wantValue: "",
			wantOk:    false,
		},
		{
			name:      "nil error",
			err:       nil,
			key:       "any",
			wantValue: "",
			wantOk:    false,
		},
		{
			name:      "type mismatch",
			err:       err,
			key:       "int",
			wantValue: "",
			wantOk:    false, // Should fail when getting int as string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test with string type
			if slices.Contains([]string{"string", "nonexistent", "any"}, tt.key) {
				gotValue, gotOk := efmt.Value[string](tt.err, tt.key)

				if tt.wantOk {
					if !gotOk {
						t.Errorf("Value() ok = false, want true")
					}

					if gotValue != tt.wantValue {
						t.Errorf("Value() = %v, want %v", gotValue, tt.wantValue)
					}
				} else if gotOk {
					t.Errorf("Value() ok = true, want false")
				}
			} else if tt.key == "int" {
				if tt.name == "type mismatch" {
					// Test type mismatch
					_, gotOk := efmt.Value[string](tt.err, tt.key)
					if gotOk {
						t.Errorf("Value() with type mismatch returned ok = true, want false")
					}
				} else {
					// Test with int type
					gotValue, gotOk := efmt.Value[int](tt.err, tt.key)

					if !gotOk {
						t.Errorf("Value() ok = false, want true")
					}

					if gotValue != tt.wantValue {
						t.Errorf("Value() = %v, want %v", gotValue, tt.wantValue)
					}
				}
			} else if tt.key == "bool" {
				// Test with bool type
				gotValue, gotOk := efmt.Value[bool](tt.err, tt.key)

				if !gotOk {
					t.Errorf("Value() ok = false, want true")
				}

				if gotValue != tt.wantValue {
					t.Errorf("Value() = %v, want %v", gotValue, tt.wantValue)
				}
			} else if tt.key == "float" {
				// Test with float64 type
				gotValue, gotOk := efmt.Value[float64](tt.err, tt.key)

				if !gotOk {
					t.Errorf("Value() ok = false, want true")
				}

				if gotValue != tt.wantValue {
					t.Errorf("Value() = %v, want %v", gotValue, tt.wantValue)
				}
			}
		})
	}
}

func TestValueRespectsShadowing(t *testing.T) {
	err := efmt.Wrap(efmt.New("inner", efmt.KV("id", 42)), "outer", efmt.KV("id", "abc"))

	if got, ok := efmt.Value[string](err, "id"); !ok || got != "abc" {
		t.Errorf("Value[string]() = %q, %v, want %q, true", got, ok, "abc")
	}

	// the inner int value is shadowed by the outer string one
	if got, ok := efmt.Value[int](err, "id"); ok {
		t.Errorf("Value[int]() = %d, true, want ok = false", got)
	}
}

func TestValuesInErrorTree(t *testing.T) {
	first := efmt.New("first", efmt.KV("first", 1), efmt.KV("common", "first"))
	second := efmt.New("second", efmt.KV("second", 2), efmt.KV("common", "second"))

	tests := []struct {
		name string
		err  error
		want map[string]any
	}{
		{
			name: "errors.Join",
			err:  efmt.Wrap(errors.Join(first, second), "top", efmt.KV("top", 0)),
			want: map[string]any{"top": 0, "first": 1, "second": 2, "common": "first"},
		},
		{
			name: "multiple %w",
			err:  fmt.Errorf("%w; %w", second, first),
			want: map[string]any{"first": 1, "second": 2, "common": "second"},
		},
		{
			name: "fmt.Errorf in the middle of a chain",
			err:  efmt.Wrap(fmt.Errorf("plain: %w", first), "top", efmt.KV("top", 0)),
			want: map[string]any{"top": 0, "first": 1, "common": "first"},
		},
		{
			name: "nested join",
			err:  errors.Join(errors.New("plain"), errors.Join(nil, second)),
			want: map[string]any{"second": 2, "common": "second"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := efmt.Values(tt.err)
			if !maps.Equal(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}
