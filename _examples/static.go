// Copyright 2014 The orujo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jroimartin/orujo"
	restlog "github.com/jroimartin/orujo/handlers/log"
)

func main() {
	s := orujo.NewServer("localhost:8080")

	logger := log.New(os.Stdout, "[HELLO] ", log.LstdFlags)
	logHandler := restlog.NewLogHandler(logger,
		"{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}")

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(".")))

	s.Route("/static/.*",
		staticHandler,
		orujo.M(logHandler),
	)

	log.Fatalln(s.ListenAndServe())
}
