// Copyright 2016 SendGrid, inc. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package rho

// AssertFunc wraps type assertions on err, and is primarily used alongside
// Scan and Underlying error chains. These functions should assign true into
// *ok to indicate that the assertion has succeeded.
type AssertFunc func(err error, ok *bool)
