// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jroimartin/gorest"
	restlog "github.com/jroimartin/gorest/handlers/log"
)

func main() {
	s := gorest.NewServer("localhost:8080")

	logger := log.New(os.Stdout, "[HELLO] ", log.LstdFlags)
	logHandler := restlog.NewLogHandler(logger,
		"{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}")

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(".")))

	s.Route("/static/.*",
		staticHandler,
		gorest.M(logHandler),
	)

	log.Fatalln(s.ListenAndServe())
}
