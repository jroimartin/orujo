// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package log implements the bult-in logging handler of gorest
*/
package log

import (
	"bytes"
	"log"
	"net/http"
	"text/template"

	"github.com/jroimartin/gorest"
)

// A LogHandler is a gorest built-in handler that provides
// logging facilities.
type LogHandler struct {
	log       *log.Logger
	tmpl      *template.Template
	mandatory bool
}

// NewLogHandler allocates and returns a new LogHandler. It
// accepts a log.Logger, that will be used to write the
// generated log records using the format specified by the
// argument fmt.
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
// struct {
//     Resp   http.ResponseWriter
//     Req    *http.Request
//     Errors []error
// }
//
// This way, the LogHandler has access to the
// http.ResponseWriter, the http.Request and a errors slice
// containing all the errors registered during the handlers
// pipe execution.
//
// E.g.:
// const logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}
// {{range  $err := .Errors}}  Err: {{$err}}
// {{end}}`
// logger := log.New(os.Stdout, "[SESSION] ", log.LstdFlags)
// logHandler := restlog.NewLogHandler(logger, logLine)
func NewLogHandler(logger *log.Logger, fmt string) *LogHandler {
	tmpl := template.Must(template.New("fmt").Parse(fmt))
	return &LogHandler{log: logger, tmpl: tmpl}
}

// SetMandatory allows to set the LogHandler as mandatory or optional.
func (h *LogHandler) SetMandatory(v bool) {
	h.mandatory = v
}

// Mandatory returns if the LogHandler is mandatory.
func (h *LogHandler) Mandatory() bool {
	return h.mandatory
}

// ServeHTTP will execute the log template when the handler is used.
func (h *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	errors := gorest.Errors(w)
	data := struct {
		Resp   http.ResponseWriter
		Req    *http.Request
		Errors []error
	}{w, r, errors}

	var out bytes.Buffer
	h.tmpl.Execute(&out, data)
	h.log.Println(out.String())
}
