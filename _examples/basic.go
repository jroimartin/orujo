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
	"github.com/jroimartin/gorest/handlers/basic"
	restlog "github.com/jroimartin/gorest/handlers/log"
)

const logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}
{{range  $err := .Errors}}  Err: {{$err}}
{{end}}`

func main() {
	s := gorest.NewServer("localhost:8080")

	logger := log.New(os.Stdout, "[HELLO] ", log.LstdFlags)
	logHandler := gorest.M(restlog.NewLogHandler(logger, logLine))

	basicHandler := gorest.M(basic.NewBasicHandler("hello", "user", "password123"))

	s.Route("/hello_auth/{name}",
		basicHandler,
		gorest.H(http.HandlerFunc(helloHandler)),
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
