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
	"github.com/jroimartin/gorest/handlers/sessions"
)

var sessionHandler *sessions.SessionHandler

const logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}
{{range  $err := .Errors}}  Err: {{$err}}
{{end}}`

func main() {
	s := gorest.NewServer("localhost:8080")

	logger := log.New(os.Stdout, "[SESSION] ", log.LstdFlags)
	logHandler := restlog.NewLogHandler(logger, logLine)

	sessionHandler = sessions.NewSessionHandler("gorest", []byte("secret"))
	sessionHandler.SetOptions(&sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	})

	s.Route("/",
		sessionHandler,
		gorest.H(http.HandlerFunc(homeHandler)),
		gorest.M(logHandler),
	)

	log.Fatalln(s.ListenAndServe())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	sessionId, err := sessionHandler.SessionId(r)
	if err != nil {
		internalServerError(w)
		gorest.RegisterError(w, err)
	}
	fmt.Fprintln(w, "SessionID:", sessionId)
}

func internalServerError(w http.ResponseWriter) {
	status := http.StatusInternalServerError
	http.Error(w, http.StatusText(status), status)
}
