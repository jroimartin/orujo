// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jroimartin/gorest"
	restlog "github.com/jroimartin/gorest/handlers/log"
)

func main() {
	s := gorest.NewServer("localhost:8080")

	logger := log.New(os.Stdout, "[HELLO] ", log.LstdFlags)
	logHandler := gorest.M(restlog.NewLogHandler(logger,
		"{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}"))

	s.RouteDefault(
		gorest.H(http.NotFoundHandler()),
		logHandler,
	)

	s.Route("/hello/{name}",
		gorest.H(http.HandlerFunc(helloHandler)),
		logHandler,
	)

	log.Fatalln(s.ListenAndServe())
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	vars := gorest.Vars(r)
	fmt.Fprintln(w, "Hello,", vars["name"])
}
