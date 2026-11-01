package efmt

import (
	"errors"
	"fmt"
)

// New creates a new formatted error with the given text and optional key-value pairs.
func New(text string, kv ...KeyValue) error {
	return &e{
		err:     errors.New(text),
		wrapped: nil,
		values:  kv,
	}
}

// Wrap wraps an existing error with additional text and optional key-value pairs.
// It preserves the original error's context while adding new information.
// Returns nil if the provided error is nil.
func Wrap(err error, text string, kv ...KeyValue) error {
	if err == nil {
		return nil
	}

	return (&e{
		wrapped: err,
		err:     fmt.Errorf("%s: %w", text, err),
		values:  make([]KeyValue, 0, 1),
	}).Add(kv)
}

// KV creates a key-value pair for attaching to errors.
// This is a convenience function to create structured data entries.
func KV(key string, value any) KeyValue {
	return KeyValue{Key: key, Value: value}
}

type e struct {
	wrapped error
	err     error
	values  []KeyValue
}

func (e *e) Error() string {
	return e.err.Error()
}

func (e *e) Unwrap() error {
	return e.wrapped
}

func (e *e) Values() []KeyValue {
	return e.values
}

func (e *e) Add(keyValues []KeyValue) *e {
	e.values = append(e.values, keyValues...)

	return e
}

// KeyValue is a single piece of structured data attached to an error.
type KeyValue struct {
	Key   string
	Value any
}
