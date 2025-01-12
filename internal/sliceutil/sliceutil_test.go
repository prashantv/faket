package sliceutil_test

import (
	"strconv"
	"testing"

	"github.com/prashantv/faket/internal/sliceutil"
	"github.com/prashantv/faket/internal/want"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		fn   func(int) string
		want []string
	}{
		{
			name: "nil",
			in:   nil,
			want: nil,
		},
		{
			name: "empty",
			in:   []int{},
			want: []string{},
		},
		{
			name: "single element",
			in:   []int{1},
			want: []string{"1"},
		},
		{
			name: "multiple elements",
			in:   []int{1, 4, 7},
			want: []string{"1", "4", "7"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sliceutil.Map(tt.in, strconv.Itoa)
			want.DeepEqual(t, "Map", got, tt.want)
		})
	}
}

func TestToSet(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want map[int]struct{}
	}{
		{
			name: "nil",
			in:   nil,
			want: map[int]struct{}{},
		},
		{
			name: "empty",
			in:   []int{},
			want: map[int]struct{}{},
		},
		{
			name: "single element",
			in:   []int{1},
			want: map[int]struct{}{
				1: {},
			},
		},
		{
			name: "unique",
			in:   []int{1, 4, 7},
			want: map[int]struct{}{
				1: {},
				4: {},
				7: {},
			},
		},
		{
			name: "duplicates",
			in:   []int{1, 4, 1, 7, 4},
			want: map[int]struct{}{
				1: {},
				4: {},
				7: {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sliceutil.ToSet(tt.in)
			want.DeepEqual(t, "ToSet", got, tt.want)
		})
	}
}
