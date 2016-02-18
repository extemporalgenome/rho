// Copyright TBD

// +build ignore

package rho

// IterFunc is used to inspect a sequence or nested chain of errors, and
// returns true to indicate that iteration should continue.
type IterFunc func(err error) bool

// ToAssert wraps an IterFunc into an AssertFunc.
func ToAssert(f IterFunc) AssertFunc {
	return func(err error, ok *bool) {
		// returned false (no more iterations requested)
		*ok = !f(err)
	}
}

// ToIter wraps an AssertFunc into an IterFunc.
func ToIter(f AssertFunc) IterFunc {
	return func(err error) (ok bool) {
		f(err, &ok)
		// assertion succeeded (no more iterations needed)
		return !ok
	}
}

type CheckFunc func(error) bool

func Check(v *bool, f CheckFunc) AssertFunc {
	return func(err error, ok *bool) {
		*v = f(err)
		*ok = *v
	}
}
