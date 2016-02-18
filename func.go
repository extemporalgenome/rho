// Copyright TBD

package rho

// AssertFunc wraps type assertions on err, and is primarily used alongside
// Scan and Underlying error chains. These functions should assign true into
// *ok to indicate that the assertion has succeeded.
type AssertFunc func(err error, ok *bool)
