package efmt

// New creates a new formatted error with the given text and optional key-value pairs.
func New(text string, kv ...KeyValue) error {
	return &e{
		msg:    text,
		cause:  nil,
		values: kv,
	}
}

// Wrap wraps an existing error with additional text and optional key-value pairs.
// It preserves the original error's context while adding new information.
// Returns nil if the provided error is nil.
func Wrap(err error, text string, kv ...KeyValue) error {
	if err == nil {
		return nil
	}

	return &e{
		msg:    text,
		cause:  err,
		values: kv,
	}
}

// KV creates a key-value pair for attaching to errors.
// This is a convenience function to create structured data entries.
func KV(key string, value any) KeyValue {
	return KeyValue{Key: key, Value: value}
}

// KeyValue is a single piece of structured data attached to an error.
type KeyValue struct {
	Key   string
	Value any
}

type e struct {
	msg    string
	cause  error
	values []KeyValue
}

func (e *e) Error() string {
	if e.cause == nil {
		return e.msg
	}

	return e.msg + ": " + e.cause.Error()
}

func (e *e) Unwrap() error {
	return e.cause
}

func (e *e) Values() []KeyValue {
	return e.values
}
