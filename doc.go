// Copyright 2014 The orujo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package orujo implements a minimalist web framework, which
has been designed to work seamlessly with the standard
net/http library. A simple hello world would look like the
following snippet:

	package main

	import (
		"fmt"
		"log"
		"net/http"

		"github.com/jroimartin/orujo"
	)

	func main() {
		s := orujo.NewServer("localhost:8080")

		s.Route("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello world!")
		}))

		log.Fatalln(s.ListenAndServe())
	}
*/
package orujo
