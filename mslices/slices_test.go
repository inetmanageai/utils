package mslices_test

import (
	"fmt"
	"testing"

	"github.com/inetmanageai/utils/mslices"
	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	type Input[T string | []string] struct {
		Slice []string
		val   T
	}
	tests := []struct {
		Name     string
		Input    interface{}
		Expected bool
	}{
		{
			Name: "string contain",
			Input: Input[string]{
				Slice: []string{"apple", "banana", "cherry"},
				val:   "banana",
			},
			Expected: true,
		},
		{
			Name: "string not contain",
			Input: Input[string]{
				Slice: []string{"apple", "banana", "cherry"},
				val:   "orange",
			},
			Expected: false,
		},
		{
			Name: "slice string contain",
			Input: Input[[]string]{
				Slice: []string{"apple", "banana", "cherry"},
				val:   []string{"banana", "cherry"},
			},
			Expected: true,
		},
		{
			Name: "slice string not contain",
			Input: Input[[]string]{
				Slice: []string{"apple", "banana", "cherry"},
				val:   []string{"banana", "orange"},
			},
			Expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// --------------- Act ---------------
			result := false
			switch input := tt.Input.(type) {
			case Input[string]:
				result = mslices.Contains(input.Slice, input.val)
			case Input[[]string]:
				result = mslices.Contains(input.Slice, input.val)
			default:
				t.Fatalf("unsupported input type %T", input)
			}

			// --------------- Assert ---------------
			assert.Equal(t, tt.Expected, result)
		})
	}
}

func TestFilter(t *testing.T) {
	type Input[T any] struct {
		v []T
		f func(T) bool
	}
	tests := []struct {
		Name     string
		Input    Input[any]
		Expected []any
	}{
		// TODO: Add test cases.
		{
			Name: "Filter even numbers",
			Input: Input[any]{
				v: []any{1, 2, 3, 4, 5},
				f: func(x any) bool { return x.(int)%2 == 0 },
			},
			Expected: []any{2, 4},
		},
		{
			Name: "Filter strings by length",
			Input: Input[any]{
				v: []any{"go", "golang", "gopher"},
				f: func(x any) bool { return len(x.(string)) > 4 },
			},
			Expected: []any{"golang", "gopher"},
		},
		{
			Name: "Filter custom struct by field value",
			Input: Input[any]{
				v: []any{
					struct {
						Name string
						Age  int
					}{"Alice", 30},
					struct {
						Name string
						Age  int
					}{"Bob", 25},
					struct {
						Name string
						Age  int
					}{"Charlie", 35},
				},
				f: func(x any) bool {
					return x.(struct {
						Name string
						Age  int
					}).Age > 30
				},
			},
			Expected: []any{
				struct {
					Name string
					Age  int
				}{"Charlie", 35},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// --------------- Act ---------------
			result := mslices.Filter(tt.Input.v, tt.Input.f)

			// --------------- Assert ---------------
			assert.Equal(t, tt.Expected, result)
		})
	}
}

func TestSome(t *testing.T) {
	type Input[T any] struct {
		v []T
		f func(T) bool
	}
	tests := []struct {
		Name     string
		Input    Input[any]
		Expected bool
	}{
		// TODO: Add test cases.
		{
			Name: "Some int, expected true",
			Input: Input[any]{
				v: []any{1, 2, 3, 4, 5},
				f: func(v any) bool {
					return v.(int) > 3
				},
			},
			Expected: true,
		},
		{
			Name: "Some int, expected false",
			Input: Input[any]{
				v: []any{1, 2, 3, 4, 5},
				f: func(v any) bool {
					return v.(int) > 5
				},
			},
			Expected: false,
		},
		{
			Name: "Some string, expected true",
			Input: Input[any]{
				v: []any{"apple", "banana", "cherry"},
				f: func(v any) bool {
					return v.(string) == "banana"
				},
			},
			Expected: true,
		},
		{
			Name: "Some string, expected false",
			Input: Input[any]{
				v: []any{"apple", "banana", "cherry"},
				f: func(v any) bool {
					return v.(string) == "orange"
				},
			},
			Expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// --------------- Act ---------------
			result := mslices.Some(tt.Input.v, tt.Input.f)

			// --------------- Assert ---------------
			assert.Equal(t, tt.Expected, result)
		})
	}
}

func TestEvery(t *testing.T) {
	type Input[T any] struct {
		v []T
		f func(T) bool
	}
	tests := []struct {
		Name     string
		Input    Input[any]
		Expected bool
	}{
		// TODO: Add test cases.
		{
			Name: "Every int, expected true",
			Input: Input[any]{
				v: []any{2, 4, 6, 8},
				f: func(v any) bool {
					return v.(int)%2 == 0
				},
			},
			Expected: true,
		},
		{
			Name: "Every int, expected false",
			Input: Input[any]{
				v: []any{2, 4, 5, 8},
				f: func(v any) bool {
					return v.(int)%2 == 0
				},
			},
			Expected: false,
		},
		{
			Name: "Every string, expected true",
			Input: Input[any]{
				v: []any{"apple", "banana", "cherry"},
				f: func(v any) bool {
					return len(v.(string)) > 4
				},
			},
			Expected: true,
		},
		{
			Name: "Every string, expected false",
			Input: Input[any]{
				v: []any{"apple", "banana", "fig"},
				f: func(v any) bool {
					return len(v.(string)) > 4
				},
			},
			Expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// --------------- Act ---------------
			result := mslices.Every(tt.Input.v, tt.Input.f)

			// --------------- Assert ---------------
			assert.Equal(t, tt.Expected, result)
		})
	}
}

func TestMap(t *testing.T) {
	type Input[T any, R any] struct {
		v []T
		f func(T) R
	}
	tests := []struct {
		Name     string
		Input    Input[any, any]
		Expected []any
	}{
		// TODO: Add test cases.
		{
			Name: "Map int to int",
			Input: Input[any, any]{
				v: []any{1, 2, 3, 4, 5},
				f: func(v any) any {
					return v.(int) * 2
				},
			},
			Expected: []any{2, 4, 6, 8, 10},
		},
		{
			Name: "Map string to string length",
			Input: Input[any, any]{
				v: []any{"apple", "banana", "cherry"},
				f: func(v any) any {
					return len(v.(string))
				},
			},
			Expected: []any{5, 6, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// --------------- Act ---------------
			result := mslices.Map(tt.Input.v, tt.Input.f)

			// --------------- Assert ---------------
			assert.Equal(t, tt.Expected, result)
		})
	}
}

func TestFind(t *testing.T) {
	type Input[T any] struct {
		v []T
		f func(T) bool
	}
	type Output[T any] struct {
		Found T
		Idx   int
		Err   error
	}
	tests := []struct {
		Name     string
		Input    Input[any]
		Expected Output[any]
	}{
		// TODO: Add test cases.
		{
			Name: "Find int, expect found",
			Input: Input[any]{
				v: []any{1, 2, 3, 4, 5},
				f: func(v any) bool {
					return v.(int) == 3
				},
			},
			Expected: Output[any]{
				Found: 3,
				Idx:   2,
				Err:   nil,
			},
		},
		{
			Name: "Find int, expect not found",
			Input: Input[any]{
				v: []any{1, 2, 3, 4, 5},
				f: func(v any) bool {
					return v.(int) == 6
				},
			},
			Expected: Output[any]{
				Found: nil,
				Idx:   -1,
				Err:   fmt.Errorf("not found"),
			},
		},
		{
			Name: "Find string, expect found",
			Input: Input[any]{
				v: []any{"apple", "banana", "cherry"},
				f: func(v any) bool {
					return v.(string) == "banana"
				},
			},
			Expected: Output[any]{
				Found: "banana",
				Idx:   1,
				Err:   nil,
			},
		},
		{
			Name: "Find string, expect not found",
			Input: Input[any]{
				v: []any{"apple", "banana", "cherry"},
				f: func(v any) bool {
					return v.(string) == "orange"
				},
			},
			Expected: Output[any]{
				Found: nil,
				Idx:   -1,
				Err:   fmt.Errorf("not found"),
			},
		},
		{
			Name: "Find custom struct, expect found",
			Input: Input[any]{
				v: []any{
					struct {
						Name string
						Age  int
					}{"Alice", 30},
					struct {
						Name string
						Age  int
					}{"Bob", 25},
					struct {
						Name string
						Age  int
					}{"Charlie", 35},
				},
				f: func(x any) bool {
					return x.(struct {
						Name string
						Age  int
					}).Age >= 30
				},
			},
			Expected: Output[any]{
				Found: struct {
					Name string
					Age  int
				}{"Alice", 30},
				Idx: 0,
				Err: nil,
			},
		},
		{
			Name: "Find custom struct, expect not found",
			Input: Input[any]{
				v: []any{
					struct {
						Name string
						Age  int
					}{"Alice", 30},
					struct {
						Name string
						Age  int
					}{"Bob", 25},
					struct {
						Name string
						Age  int
					}{"Charlie", 35},
				},
				f: func(x any) bool {
					return x.(struct {
						Name string
						Age  int
					}).Age > 40
				},
			},
			Expected: Output[any]{
				Found: nil,
				Idx:   -1,
				Err:   fmt.Errorf("not found"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// --------------- Act ---------------
			found, idx, err := mslices.Find(tt.Input.v, tt.Input.f)

			// --------------- Assert ---------------
			assert.Equal(t, tt.Expected, Output[any]{found, idx, err})
		})
	}
}

func TestSetUnique(t *testing.T) {
	tests := []struct {
		Name     string
		Input    []any
		Expected []any
	}{
		// TODO: Add test cases.
		{
			Name:     "SetUnique int",
			Input:    []any{1, 2, 2, 3, 4, 4, 5},
			Expected: []any{1, 2, 3, 4, 5},
		},
		{
			Name:     "SetUnique string",
			Input:    []any{"apple", "banana", "apple", "cherry", "banana"},
			Expected: []any{"apple", "banana", "cherry"},
		},
		{
			Name:     "SetUnique empty slice",
			Input:    []any{},
			Expected: []any{},
		},
		{
			Name:     "SetUnique single element",
			Input:    []any{1},
			Expected: []any{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// --------------- Act ---------------
			result := mslices.SetUnique(tt.Input)

			// --------------- Assert ---------------
			assert.Equal(t, tt.Expected, result)
		})
	}
}
