package server

import "net/http"

type pipe struct {
	handlers []Handler
	quit     bool
}

type pipeWriter struct {
	p *pipe
	http.ResponseWriter
}

func newPipe(handlers ...Handler) *pipe {
	return &pipe{handlers: handlers}
}

func (p *pipe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func newPipeWriter(p *pipe, w http.ResponseWriter) *pipeWriter {
	return &pipeWriter{p: p, ResponseWriter: w}
}

func (pw *pipeWriter) WriteHeader(code int) {
	pw.p.quit = true
	pw.ResponseWriter.WriteHeader(code)
}
