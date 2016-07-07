// Copyright 2016 SendGrid, inc. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package rho provides nested-error management and a variety of
// error-handling functions.
//
// Errors should use the Error method to provide a human-readable message
// and, if needed, expose error data through additional methods. TODO
//
// The Underlying method should not be used to establish chronology or
// relationships between distinct errors (parent/child or multi-errors), or
// in any way that masks contextual methods (Temporary, Status, etc). An
// Underlying chain may be loosely ordered, will often end in a regular
// error (or an Underlying method that returns nil), and each Error() method
// in an Underlying chain should provide a meaningful human readable message
// for the overall error.
package rho
