package efmt

import "log/slog"

// Attrs returns the values attached to err and the errors it wraps as slog attributes.
// Keys are deduplicated with the same precedence as Values, and attributes are
// ordered by the first occurrence of their key in ValuePairs, so the output is stable.
func Attrs(err error) []slog.Attr {
	pairs := ValuePairs(err)
	if len(pairs) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(pairs))
	res := make([]slog.Attr, 0, len(pairs))

	for _, kv := range pairs {
		if _, ok := seen[kv.Key]; ok {
			continue
		}

		seen[kv.Key] = struct{}{}

		res = append(res, slog.Any(kv.Key, kv.Value))
	}

	return res
}
