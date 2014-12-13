// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gorest

import "net/http"

type pipeHandler struct {
	handler   http.Handler
	mandatory bool
}

func (h pipeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(w, r)
}

type pipeContext struct {
	errors []error
	quit   bool
}

func newPipeContext() *pipeContext {
	ctx := &pipeContext{}
	ctx.quit = false
	ctx.errors = make([]error, 0)
	return ctx
}

type pipe struct {
	handlers []pipeHandler
}

func newPipe(handlers ...http.Handler) *pipe {
	p := &pipe{}
	p.handlers = make([]pipeHandler, 0)
	for _, h := range handlers {
		var ph pipeHandler
		ph, isPipeHandler := h.(pipeHandler)
		if !isPipeHandler {
			ph = pipeHandler{handler: h, mandatory: false}
		}
		p.handlers = append(p.handlers, ph)
	}
	return p
}

func (p *pipe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newPipeContext()

	for _, ph := range p.handlers {
		if ctx.quit && !ph.mandatory {
			continue
		}
		pw := newPipeWriter(ctx, w)
		ph.handler.ServeHTTP(pw, r)
	}
}

type pipeWriter struct {
	ctx *pipeContext
	http.ResponseWriter
}

func newPipeWriter(ctx *pipeContext, w http.ResponseWriter) *pipeWriter {
	return &pipeWriter{ctx: ctx, ResponseWriter: w}
}

func (pw *pipeWriter) WriteHeader(code int) {
	pw.ctx.quit = true
	pw.ResponseWriter.WriteHeader(code)
}

// M is a helper to set a handler as "mandatory".
func M(h http.Handler) http.Handler {
	return pipeHandler{handler: h, mandatory: true}
}

// RegisterError can be used by Handlers to register errors.
func RegisterError(w http.ResponseWriter, err error) {
	pw, isPipeWriter := w.(*pipeWriter)
	if !isPipeWriter || err == nil {
		return
	}
	pw.ctx.errors = append(pw.ctx.errors, err)
}

// Errors is used to retrieve the errors registered via RegisterError()
// during the execution of the handlers pipe.
func Errors(w http.ResponseWriter) []error {
	pw, isPipeWriter := w.(*pipeWriter)
	if isPipeWriter {
		return pw.ctx.errors
	}
	return nil
}
