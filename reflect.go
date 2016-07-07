// Copyright 2016 SendGrid, inc. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// +build ignore

package rho

import "reflect"

type pair struct {
	v reflect.Value
	t reflect.Type
}

func ScanVals(err error, a ...interface{}) bool {
	if err == nil {
		return false
	}
	n := len(a)
	// pairs stores precomputed reflect Values and Types
	pairs := make([]pair, n)
	for i := range a {
		v := reflect.ValueOf(a[i]).Elem()
		pairs[i] = pair{v, v.Type()}
	}
	for n > 0 && err != nil {
		errv := reflect.ValueOf(err)
		errt := errv.Type()
		for i := 0; i < n; {
			p := pairs[i]
			if !errt.AssignableTo(p.t) {
				i++
				continue
			}
			p.v.Set(errv)
			n--
			pairs[i], pairs[n] = pairs[n], pairs[i]
		}
		if n <= 0 {
			return true
		}
		err = Sub(err)
	}
	return false
}
