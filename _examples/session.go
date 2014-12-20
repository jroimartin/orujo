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
	restlog "github.com/jroimartin/orujo/handlers/log"
	"github.com/jroimartin/orujo/handlers/sessions"
)

var sessionHandler sessions.SessionHandler

const logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}
{{range  $err := .Errors}}  Err: {{$err}}
{{end}}`

func main() {
	s := orujo.NewServer("localhost:8080")

	logger := log.New(os.Stdout, "[SESSION] ", log.LstdFlags)
	logHandler := restlog.NewLogHandler(logger, logLine)

	sessionHandler = sessions.NewSessionHandler("orujo", []byte("s3cre7"))
	sessionHandler.SetOptions(&sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	})

	s.Route("/",
		orujo.M(sessionHandler),
		http.HandlerFunc(homeHandler),
		orujo.M(logHandler),
	)

	log.Fatalln(s.ListenAndServe())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	sessionId, err := sessionHandler.SessionID(r)
	if err != nil {
		internalServerError(w)
		orujo.RegisterError(w, err)
	}
	fmt.Fprintln(w, "SessionID:", sessionId)
}

func internalServerError(w http.ResponseWriter) {
	status := http.StatusInternalServerError
	http.Error(w, http.StatusText(status), status)
}
