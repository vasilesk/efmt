package efmt_test

import (
	"errors"
	"log/slog"
	"slices"
	"testing"

	"github.com/vasilesk/efmt"
)

func TestAttrs(t *testing.T) {
	base := efmt.New("base", efmt.KV("base", "value"), efmt.KV("common", "base"))
	top := efmt.Wrap(base, "top", efmt.KV("top", true), efmt.KV("common", "top"))

	tests := []struct {
		name string
		err  error
		want []slog.Attr
	}{
		{
			name: "nil error",
			err:  nil,
			want: nil,
		},
		{
			name: "standard error",
			err:  errors.New("standard"),
			want: nil,
		},
		{
			name: "chain with shadowed key",
			err:  top,
			want: []slog.Attr{
				slog.Bool("top", true),
				slog.String("common", "top"),
				slog.String("base", "value"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := efmt.Attrs(tt.err)
			if !slices.EqualFunc(got, tt.want, slog.Attr.Equal) {
				t.Errorf("Attrs() = %v, want %v", got, tt.want)
			}
		})
	}
}
