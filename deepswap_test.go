package deepswap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeepswap(t *testing.T) {
	tt := []struct {
		src  any
		old  any
		new  any
		resF func(any) bool
	}{
		{
			src: 1,
			old: 1,
			new: 2,
			resF: func(x any) bool {
				return x.(int) == 2
			},
		},
		{
			src: toPointer(1),
			old: 1,
			new: 2,
			resF: func(x any) bool {
				return *(x.(*int)) == 2
			},
		},
		{
			src: toPointer[*int](nil),
			old: 1,
			new: 2,
			resF: func(x any) bool {
				return *(x.(**int)) == nil
			},
		},
		{
			src: toPointer(toPointer(1)),
			old: 1,
			new: 2,
			resF: func(x any) bool {
				return **(x.(**int)) == 2
			},
		},
		{
			src: toInterface(1),
			old: 1,
			new: 2,
			resF: func(x any) bool {
				return x.(int) == 2
			},
		},
		{
			src: toInterface(toPointer(1)),
			old: 1,
			new: 2,
			resF: func(x any) bool {
				return *x.(*int) == 2
			},
		},
		{
			src: map[int]int{
				1: 1,
				2: 2,
				4: 1,
			},
			old: 1,
			new: 2,
			resF: func(x any) bool {
				m := x.(map[int]int)
				return m[1] == 2 && m[2] == 2 && m[4] == 2
			},
		},
		{
			src: []int{1, 1, 1},
			old: 1,
			new: 2,
			resF: func(x any) bool {
				s := x.([]int)
				return s[0] == 2 && s[1] == 2 && s[2] == 2
			},
		},
		{
			src: [...]int{1, 1, 1},
			old: 1,
			new: 2,
			resF: func(x any) bool {
				s := x.([3]int)
				return s[0] == 2 && s[1] == 2 && s[2] == 2
			},
		},
		{
			src: struct{ X int }{X: 1},
			old: 1,
			new: 2,
			resF: func(x any) bool {
				s := x.(struct{ X int })
				return s.X == 2
			},
		},
		{
			src: struct{ X *int }{X: toPointer(1)},
			old: 1,
			new: 2,
			resF: func(x any) bool {
				s := x.(struct{ X *int })
				return *s.X == 2
			},
		},
		{
			src: toInterface(struct{ X []int }{
				X: []int{1, 1},
			}),
			old: 1,
			new: 2,
			resF: func(x any) bool {
				s := x.(struct{ X []int })
				return s.X[0] == 2 && s.X[1] == 2
			},
		},
		{
			src: struct {
				A string
				B int
				C *int
				D []int
				E map[string]string
			}{
				A: "x",
				B: 1,
				C: toPointer(1),
				D: []int{1, 1},
				E: nil,
			},
			old: 1,
			new: 2,
			resF: func(x any) bool {
				s := x.(struct {
					A string
					B int
					C *int
					D []int
					E map[string]string
				})
				return s.A == "x" && s.B == 2 && *s.C == 2 && s.D[0] == 2 && s.D[1] == 2 && s.E == nil
			},
		},
	}

	for _, tc := range tt {
		res := DeepSwap(tc.src, tc.old, tc.new)

		assert.Truef(t, tc.resF(res), "[%d] %v does not satisfy condition ", res)
	}
}
