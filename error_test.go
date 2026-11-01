package efmt_test

import (
	"errors"
	"testing"

	"github.com/vasilesk/efmt"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		kv      []efmt.KeyValue
		wantErr string
	}{
		{
			name:    "simple error",
			text:    "simple error",
			kv:      nil,
			wantErr: "simple error",
		},
		{
			name:    "error with key-value",
			text:    "error with context",
			kv:      []efmt.KeyValue{efmt.KV("key", "value")},
			wantErr: "error with context",
		},
		{
			name:    "error with multiple key-values",
			text:    "complex error",
			kv:      []efmt.KeyValue{efmt.KV("key1", "value1"), efmt.KV("key2", 42)},
			wantErr: "complex error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := efmt.New(tt.text, tt.kv...)
			if err.Error() != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err.Error(), tt.wantErr)

				return
			}

			// Check values
			if e, ok := err.(interface{ Values() []efmt.KeyValue }); ok {
				if len(e.Values()) != len(tt.kv) {
					t.Errorf("New() values count = %v, want %v", len(e.Values()), len(tt.kv))
				}
			} else {
				t.Errorf("New() error doesn't implement Values() method")
			}
		})
	}
}

func TestWrap(t *testing.T) {
	baseErr := errors.New("base error")

	tests := []struct {
		name    string
		err     error
		text    string
		kv      []efmt.KeyValue
		wantErr string
		wantNil bool
	}{
		{
			name:    "wrap nil error",
			err:     nil,
			text:    "wrapper",
			kv:      nil,
			wantNil: true,
		},
		{
			name:    "wrap simple error",
			err:     baseErr,
			text:    "wrapper",
			kv:      nil,
			wantErr: "wrapper: base error",
		},
		{
			name:    "wrap with key-value",
			err:     baseErr,
			text:    "wrapper",
			kv:      []efmt.KeyValue{efmt.KV("key", "value")},
			wantErr: "wrapper: base error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := efmt.Wrap(tt.err, tt.text, tt.kv...)

			if tt.wantNil {
				if err != nil {
					t.Errorf("Wrap() error = %v, want nil", err)
				}

				return
			}

			if err.Error() != tt.wantErr {
				t.Errorf("Wrap() error = %v, wantErr %v", err.Error(), tt.wantErr)

				return
			}

			// Check unwrapping
			unwrapped := errors.Unwrap(err)
			if unwrapped != tt.err { //nolint:errorlint // the direct cause itself is expected here
				t.Errorf("Unwrap() = %v, want %v", unwrapped, tt.err)
			}

			// Check values
			if e, ok := err.(interface{ Values() []efmt.KeyValue }); ok {
				if len(e.Values()) != len(tt.kv) {
					t.Errorf("Wrap() values count = %v, want %v", len(e.Values()), len(tt.kv))
				}
			} else {
				t.Errorf("Wrap() error doesn't implement Values() method")
			}
		})
	}
}

func TestKV(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value any
	}{
		{
			name:  "string value",
			key:   "key",
			value: "value",
		},
		{
			name:  "int value",
			key:   "count",
			value: 42,
		},
		{
			name:  "bool value",
			key:   "valid",
			value: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := efmt.KV(tt.key, tt.value)
			if kv.Key != tt.key {
				t.Errorf("KV().Key = %v, want %v", kv.Key, tt.key)
			}

			if kv.Value != tt.value {
				t.Errorf("KV().Value = %v, want %v", kv.Value, tt.value)
			}
		})
	}
}

func TestKeyValuesAreCopied(t *testing.T) {
	kvs := []efmt.KeyValue{efmt.KV("key", "original")}

	newErr := efmt.New("new", kvs...)
	wrapErr := efmt.Wrap(errors.New("base"), "wrap", kvs...)

	kvs[0].Value = "mutated"

	for name, err := range map[string]error{"New": newErr, "Wrap": wrapErr} {
		if got, _ := efmt.Value[string](err, "key"); got != "original" {
			t.Errorf("%s() value changed after caller mutated its slice: got %q", name, got)
		}
	}
}

type codeError struct{ code int }

func (c *codeError) Error() string { return "code error" }

func TestWrapIsAndAs(t *testing.T) {
	sentinel := errors.New("sentinel")
	coded := &codeError{code: 404}

	err := efmt.Wrap(efmt.Wrap(errors.Join(sentinel, coded), "inner"), "outer", efmt.KV("key", "value"))

	if !errors.Is(err, sentinel) {
		t.Errorf("errors.Is() = false, want true")
	}

	var target *codeError
	if !errors.As(err, &target) || target.code != 404 {
		t.Errorf("errors.As() did not find *codeError")
	}

	if !errors.Is(efmt.Wrap(sentinel, "wrapped"), sentinel) {
		t.Errorf("errors.Is() through a single Wrap = false, want true")
	}
}
