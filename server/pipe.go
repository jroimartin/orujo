// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import "net/http"

type pipe struct {
	handlers []Handler
	quit     bool
}

func newPipe(handlers ...Handler) *pipe {
	return &pipe{handlers: handlers}
}

func (p *pipe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.quit = false

	for _, h := range p.handlers {
		if h == nil {
			continue
		}
		if p.quit && !h.Mandatory() {
			continue
		}
		pw := newPipeWriter(p, w)
		h.ServeHTTP(pw, r)
	}
}

type pipeWriter struct {
	p *pipe
	http.ResponseWriter
}

func newPipeWriter(p *pipe, w http.ResponseWriter) *pipeWriter {
	return &pipeWriter{p: p, ResponseWriter: w}
}

func (pw *pipeWriter) WriteHeader(code int) {
	pw.p.quit = true
	pw.ResponseWriter.WriteHeader(code)
}
