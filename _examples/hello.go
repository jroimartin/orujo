// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	restlog "github.com/jroimartin/gorest/handlers/log"
	"github.com/jroimartin/gorest/server"
)

func main() {
	s := server.NewServer("localhost:8080")

	logger := log.New(os.Stdout, "[HELLO] ", log.LstdFlags)
	logHandler := restlog.NewLogHandler(logger,
		"{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}")

	s.Route("/hello/{name}",
		server.H(http.HandlerFunc(helloHandler)),
		server.M(logHandler),
	)

	log.Fatalln(s.ListenAndServe())
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintln(w, "Hello,", vars["name"])
}
