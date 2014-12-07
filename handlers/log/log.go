// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"bytes"
	"log"
	"net/http"
	"text/template"

	"github.com/jroimartin/gorest/server"
)

type LogHandler struct {
	log       *log.Logger
	tmpl      *template.Template
	mandatory bool
}

func NewLogHandler(logger *log.Logger, fmt string) *LogHandler {
	tmpl := template.Must(template.New("fmt").Parse(fmt))
	return &LogHandler{log: logger, tmpl: tmpl}
}

func (h *LogHandler) SetMandatory(v bool) {
	h.mandatory = v
}

func (h *LogHandler) Mandatory() bool {
	return h.mandatory
}

func (h *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var errors []error

	pw, isPipeWriter := w.(*server.PipeWriter)
	if isPipeWriter {
		errors = pw.Errors()
	}

	data := struct {
		Resp   http.ResponseWriter
		Req    *http.Request
		Errors []error
	}{w, r, errors}

	var out bytes.Buffer
	h.tmpl.Execute(&out, data)
	h.log.Println(out.String())
}