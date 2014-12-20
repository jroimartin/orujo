// Copyright 2014 The orujo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orujo

import "net/http"

// A Server represents an HTTP server that will be used to serve the
// web service.
type Server struct {
	// Addr allows to configure the TCP network address used by the
	// web service. It must be set before calling ListenAndServe or
	// ListenAndServeTLS.
	Addr string

	mux    *http.ServeMux
	router *router
}

// NewServer returns a new Sever that will listen on addr.
func NewServer(addr string) Server {
	var s Server
	s.Addr = addr
	s.mux = http.NewServeMux()
	s.router = newRouter()
	s.mux.Handle("/", s.router)
	return s
}

// ListenAndServe listens on the TCP network address s.Addr and then
// handles requests on incoming connections.
func (s Server) ListenAndServe() error {
	return http.ListenAndServe(s.Addr, s.mux)
}

// ListenAndServeTLS acts identically to ListenAndServe, except that it
// expects HTTPS connections. For more information see the documentation
// of the package net/http.
func (s Server) ListenAndServeTLS(certFile string, keyFile string) error {
	return http.ListenAndServeTLS(s.Addr, certFile, keyFile, s.mux)
}

// Route registers a new route for the specified URL path and allows to
// register the handlers pipe that will be used to handle the requests
// to that resource.
//
// path is a regular expression and will panic if it cannot be compiled.
// See the documentation of the package regexp.
//
// The handlers pipe is executed sequentially until a HTTP header is
// explicitally written in the response. Besides that, mandatory handlers
// are always executed.
func (s Server) Route(path string, handlers ...http.Handler) *Route {
	return s.router.handle(path, newPipe(handlers...))
}

// RouteDefault sets the default route. This handler is used if any other
// route matches the request URI.
func (s Server) RouteDefault(handlers ...http.Handler) {
	s.router.defaultRoute = newPipe(handlers...)
}
