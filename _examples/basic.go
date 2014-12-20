// Copyright 2014 The orujo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jroimartin/orujo"
	"github.com/jroimartin/orujo/handlers/basic"
	restlog "github.com/jroimartin/orujo/handlers/log"
)

const logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}
{{range  $err := .Errors}}  Err: {{$err}}
{{end}}`

func main() {
	s := orujo.NewServer("localhost:8080")

	logger := log.New(os.Stdout, "[HELLO] ", log.LstdFlags)
	logHandler := restlog.NewLogHandler(logger, logLine)

	basicHandler := basic.NewBasicHandler("hello", "user", "password123")

	s.Route(`/hello/\w+`,
		basicHandler,
		http.HandlerFunc(helloHandler),
		orujo.M(logHandler),
	)

	log.Fatalln(s.ListenAndServe())
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	vars := orujo.Vars(r, `/hello/(?P<name>\w+)`)
	fmt.Fprintln(w, "Hello,", vars["name"])
}
