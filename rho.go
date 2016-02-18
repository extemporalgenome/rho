// Copyright TBD

package rho

import "fmt"

// Error is implemented by values that can wrap a nested error. Underlying
// may return nil if there is no nested error.
type Error interface {
	error
	Underlying() error
}

// Errorf is a drop-in replacement for fmt.Errorf, but additionally exposes
// an Underlying method which will return the first error passed in, or nil
// if none of the arguments implement error.
func Errorf(format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	for _, v := range a {
		err, ok := v.(error)
		if ok {
			return &nested{msg, err}
		}
	}
	return &nested{msg, nil}
}

type nested struct {
	msg string
	err error
}

func (n *nested) Error() string     { return n.msg }
func (n *nested) Underlying() error { return n.err }

// Sub returns the nested error if err also implements Error, or nil
// otherwise. An error which does not have an Underlying method is
// considered equivalent to an Underlying method that returns nil.
func Sub(err error) error {
	e, ok := err.(Error)
	if ok {
		return e.Underlying()
	}
	return nil
}

// Scan recurses through a chain of errors until a nil nested error is
// encountered. Each AssertFunc will be called at every level of recursion
// (including the initial error) until it assigns true into its *bool
// parameter. These callbacks may be invoked in arbitrary order, and the
// variadic slice order may be rearranged. Scan returns true if and only if
// every callback succeeded.
//
// If the outermost error is nil, none of the supplied functions will be
// called. Scan follows the nested-error semantics of Sub.
func Scan(err error, funcs ...AssertFunc) bool {
	var ok bool
	n := len(funcs)
	for n > 0 && err != nil {
		for i := 0; i < n; {
			ok = false
			funcs[i](err, &ok)
			if !ok {
				i++
				continue
			}
			n--
			// func no longer needed; swap to end of slice
			// (avoid n^2 behavior this way).
			funcs[i], funcs[n] = funcs[n], funcs[i]
		}
		if n <= 0 {
			// all functions succeeded
			return true
		}
		err = Sub(err)
	}
	return false
}
