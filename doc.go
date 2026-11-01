// Package efmt attaches structured key-value data to errors.
//
// Errors are created with [New] or wrapped with [Wrap], passing pairs built by [KV].
// The data can later be read back from anywhere in the error tree with
// [Values], [Value], [ValuePairs] or, for logging, [Attrs].
//
// Errors produced by this package work with [errors.Is], [errors.As] and
// [errors.Unwrap]. Values are also collected through errors that were wrapped
// by other means, such as fmt.Errorf with %w or [errors.Join].
//
// When the same key is set at several levels, the outermost value wins:
// wrapping an error can override what was attached closer to its origin.
package efmt
