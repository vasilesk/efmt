package efmt

import "errors"

// ValuePairs extracts all key-value pairs from an error and its wrapped errors.
// It traverses the error chain and collects structured data from each error.
func ValuePairs(err error) []KeyValue {
	var res []KeyValue

	for err := err; err != nil; {
		if e, ok := err.(*e); ok {
			res = append(res, e.Values()...)
		}

		err = errors.Unwrap(err)
	}

	return res
}

// Values extracts all key-value pairs from an error chain and converts them to a map.
// When multiple errors contain the same key, the value from the outermost error is used.
func Values(err error) map[string]any {
	pairs := ValuePairs(err)
	res := make(map[string]any, len(pairs))

	// pairs are ordered from the outermost error, so the first occurrence wins
	for _, kv := range pairs {
		if _, ok := res[kv.Key]; !ok {
			res[kv.Key] = kv.Value
		}
	}

	return res
}

// Value retrieves a typed value for a specific key from the error chain.
// As with Values, the outermost occurrence of the key is used.
// It returns false if the key is not found or if that value is not of type T;
// in the latter case inner values with the same key are not consulted.
func Value[T any](err error, key string) (T, bool) {
	for _, kv := range ValuePairs(err) {
		if kv.Key == key {
			val, ok := kv.Value.(T)

			return val, ok
		}
	}

	var zero T

	return zero, false
}
