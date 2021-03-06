// Copyright 2014 The orujo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package log implements the bult-in logging handler of orujo.
*/
package log

import (
	"bytes"
	"log"
	"net/http"
	"text/template"

	"github.com/jroimartin/orujo"
)

// A LogHandler is a orujo built-in handler that provides
// logging features.
type LogHandler struct {
	log  *log.Logger
	tmpl *template.Template
}

// NewLogHandler returns a new LogHandler. It accepts a
// log.Logger, that will be used to write the generated
// log records using the format specified by the argument
// fmt.
//
// fmt must be a valid text template, which is parsed when
// NewLogHandler is called. Note that if fmt is not valid,
// NewLogHandler will panic. For more information regarding
// templates, see the documentation of the package
// "text/template".
//
// The template will be executed passing the following
// structure:
//
//     struct {
//         Resp   http.ResponseWriter
//         Req    *http.Request
//         Errors []error
//     }
//
// This way, the LogHandler has access to the
// http.ResponseWriter, the http.Request and a errors slice
// containing all the errors registered during the handlers
// pipe execution.
//
// E.g.:
//     const logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}
//     {{range  $err := .Errors}}  Err: {{$err}}
//     {{end}}`
//     logger := log.New(os.Stdout, "[ORUJO] ", log.LstdFlags)
//     logHandler := orujolog.NewLogHandler(logger, logLine)
func NewLogHandler(logger *log.Logger, fmt string) LogHandler {
	tmpl := template.Must(template.New("fmt").Parse(fmt))
	return LogHandler{log: logger, tmpl: tmpl}
}

// ServeHTTP will execute the log template when the handler is used.
func (h LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	errors := orujo.Errors(w)
	data := struct {
		Resp   http.ResponseWriter
		Req    *http.Request
		Errors []error
	}{w, r, errors}

	var out bytes.Buffer
	h.tmpl.Execute(&out, data)
	h.log.Println(out.String())
}
