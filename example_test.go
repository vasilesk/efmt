package efmt_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/vasilesk/efmt"
)

var errNotFound = errors.New("not found")

func findUser(id int) error {
	return efmt.Wrap(errNotFound, "user lookup failed", efmt.KV("user_id", id))
}

func ExampleNew() {
	err := efmt.New("user is underage", efmt.KV("age", 16), efmt.KV("required_age", 18))

	fmt.Println(err)
	fmt.Println(efmt.Values(err))
	// Output:
	// user is underage
	// map[age:16 required_age:18]
}

func ExampleWrap() {
	err := efmt.Wrap(findUser(42), "request failed", efmt.KV("request_id", "abc"))

	fmt.Println(err)
	fmt.Println(errors.Is(err, errNotFound))
	fmt.Println(efmt.Values(err))
	// Output:
	// request failed: user lookup failed: not found
	// true
	// map[request_id:abc user_id:42]
}

func ExampleValue() {
	err := findUser(42)

	if id, ok := efmt.Value[int](err, "user_id"); ok {
		fmt.Println("user id:", id)
	}

	if _, ok := efmt.Value[string](err, "user_id"); !ok {
		fmt.Println("user id is not a string")
	}
	// Output:
	// user id: 42
	// user id is not a string
}

func ExampleValues_override() {
	inner := efmt.New("inner", efmt.KV("stage", "parse"), efmt.KV("line", 7))
	outer := efmt.Wrap(inner, "outer", efmt.KV("stage", "load"))

	fmt.Println(efmt.Values(outer))
	// Output:
	// map[line:7 stage:load]
}

func ExampleValues_join() {
	err := errors.Join(
		efmt.New("first", efmt.KV("first", 1)),
		fmt.Errorf("second: %w", efmt.New("inner", efmt.KV("second", 2))),
	)

	fmt.Println(efmt.Values(err))
	// Output:
	// map[first:1 second:2]
}

func ExampleAttrs() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}

			return a
		},
	}))

	err := efmt.Wrap(findUser(42), "request failed", efmt.KV("request_id", "abc"))

	attrs := append(efmt.Attrs(err), slog.Any("error", err))
	logger.LogAttrs(context.Background(), slog.LevelError, "request failed", attrs...)
	// Output:
	// level=ERROR msg="request failed" request_id=abc user_id=42 error="request failed: user lookup failed: not found"
}
