// Copyright TBD

package rho

import (
	"errors"
	"testing"
)

func TestSub(t *testing.T) {
	tests := []struct {
		msg string
		err error
	}{
		{"", nil},
		{"", errors.New("err1")},
		{"", &nested{"err1", nil}},
		{"err1", &nested{"err2", errors.New("err1")}},
		{"err1", &nested{"err2", &nested{"err1", nil}}},
	}
	for _, test := range tests {
		msg := ""
		err := Sub(test.err)
		if err != nil {
			msg = err.Error()
		}
		if msg != test.msg {
			t.Errorf("Sub(%v) = %q, want %q", test.err, msg, test.msg)
		}
	}
}

func bound(depth int) error {
	if depth == 0 {
		return nil
	}
	return bounded{depth, new(int)}
}

func boundcount(err error) int {
	if err == nil {
		return 0
	}
	return *err.(bounded).count
}

type bounded struct {
	depth int
	count *int
}

func (b bounded) Error() string {
	return "<bounded error>"
}

func (b bounded) Underlying() error {
	*b.count++
	if b.depth <= 1 {
		return nil
	}
	return bounded{b.depth - 1, b.count}
}

func counter(n *int, stop bool) AssertFunc {
	return func(err error, ok *bool) {
		*ok = stop
		*n++
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func TestScan(t *testing.T) {
	tests := [][]int{
		// {depth, underlying-calls, assert-calls...}
		{0, 0},             // nil error, no asserts
		{1, 0},             // depth 1 error, no asserts
		{0, 0, +0},         // nil error, uncalled assert
		{0, 0, +0, +0},     // nil error, two uncalled asserts
		{1, 0, -1},         // depth 1 error, immediate-exit assert
		{1, 0, -1, -1},     // depth 1 error, two immediate-exit asserts
		{1, 0, -1, -1, -1}, // depth 1 error, three immediate-exit asserts
		{2, 0, -1},         // depth 2 error, immediate-exit assert
		{1, 1, +1},         // depth 1 error, one assert
		{2, 2, +2},         // depth 2 error, one assert
		{2, 2, -1, +2, +2}, // depth 2 error, exit-in-first
		{2, 2, +2, -1, +2}, // depth 2 error, exit-in-middle
		{2, 2, +2, +2, -1}, // depth 2 error, exit-in-last
		{2, 2, -1, -1, +2}, // depth 2 error, cont-in-last
	}
	for tidx, test := range tests {
		depth, calls, asserts := test[0], test[1], test[2:]
		err := bound(depth)
		n := len(asserts)
		expected := make([]int, n)
		occurred := make([]int, n)
		funcs := make([]AssertFunc, n)
		for i, assert := range asserts {
			expected[i] = abs(assert)
			funcs[i] = counter(&occurred[i], assert < 0)
		}
		Scan(err, funcs...)
		if calls != boundcount(err) {
			t.Errorf("%v: got %v Underlying calls, want %v",
				tidx, boundcount(err), calls)
		}
		for i := range expected {
			if expected[i] != occurred[i] {
				t.Errorf("%v: got %v assert calls, want %v",
					tidx, occurred[i], expected[i])
			}
		}
	}
}
