// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	restlog "github.com/jroimartin/gorest/handlers/log"
	"github.com/jroimartin/gorest/handlers/sessions"
	"github.com/jroimartin/gorest/server"
)

var sessionHandler *sessions.SessionHandler

func main() {
	s := server.NewServer("localhost:8080")

	logger := log.New(os.Stdout, "[SESSION] ", log.LstdFlags)
	logHandler := server.M(restlog.NewLogHandler(logger,
		"{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}"))
	sessionHandler = sessions.NewSessionHandler("gorest", []byte("secret"))

	s.Route("/",
		sessionHandler,
		server.H(http.HandlerFunc(homeHandler)),
		logHandler,
	)

	log.Fatalln(s.ListenAndServe())
}

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	sessionId, err := sessionHandler.SessionId(r)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(rw, "SessionID:", sessionId)
}
