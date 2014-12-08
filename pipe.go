// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gorest

import "net/http"

type pipe struct {
	handlers []Handler
}

func newPipe(handlers ...Handler) *pipe {
	return &pipe{handlers: handlers}
}

func (p *pipe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newPipeContext()

	for _, h := range p.handlers {
		if h == nil {
			continue
		}
		if ctx.quit && !h.Mandatory() {
			continue
		}
		pw := newPipeWriter(ctx, w)
		h.ServeHTTP(pw, r)
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

type pipeContext struct {
	errors   []error
	quit     bool
}

func newPipeContext() *pipeContext {
	ctx := &pipeContext{}
	ctx.quit = false
	ctx.errors = make([]error, 0)
	return ctx
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
