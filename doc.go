// Copyright 2014 The orujo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package orujo solves a common problem, which is the execution of several
middlewares per route. It has been designed to work seamlessly with the
standard net/http library. A simple hello world would look like the following
snippet:

	package main

	import (
		"fmt"
		"log"
		"net/http"

		"github.com/jroimartin/orujo"
	)

	func main() {
		http.Handle("/hello", orujo.NewPipe(
			http.HandlerFunc(helloHandler),
			orujo.M(http.HandlerFunc(worldHandler)),
		))
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

	func helloHandler(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Hello, ")
	}

	func worldHandler(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "world")
	}
*/
package orujo
