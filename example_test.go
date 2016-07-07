// Copyright 2016 SendGrid, inc. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package rho_test

import (
	"errors"
	"fmt"
	"net"

	"github.com/xtgo/rho"
)

func checkDNS(name string) error {
	return &net.DNSError{
		Err:       "server timed out",
		Name:      name,
		IsTimeout: true,
	}
}

type (
	timeout interface {
		Timeout() bool
	}
	temporary interface {
		Temporary() bool
	}
)

func ExampleErrorf() {
	err1 := errors.New("original error")
	err2 := rho.Errorf("contextual failure: %v", err1)
	fmt.Println("err2:", err2)
	fmt.Println("err1 properly nested:", err1 == rho.Sub(err2))
	// Output:
	// err2: contextual failure: original error
	// err1 properly nested: true
}

func ExampleScan() {
	err := checkDNS("xyz")
	if err != nil {
		var tmo timeout
		var tmp temporary
		rho.Scan(err,
			func(err error, ok *bool) { tmo, *ok = err.(timeout) },
			func(err error, ok *bool) { tmp, *ok = err.(temporary) },
		)
		fmt.Println(err)
		if tmo != nil && tmo.Timeout() {
			fmt.Println("- operation timed out")
		}
		if tmp != nil && tmp.Temporary() {
			fmt.Println("- temporary error")
		}
	}
	// Output:
	// lookup xyz: server timed out
	// - operation timed out
	// - temporary error
}
