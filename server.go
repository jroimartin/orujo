// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gorest

import (
	"net/http"

	"github.com/gorilla/mux"
)

// A Server represents an HTTP server that will be used to serve the
// REST service.
type Server struct {
	// Addr allows to configure the TCP network address used by the
	// REST service. It must be set before calling ListenAndServe or
	// ListenAndServeTLS.
	Addr string

	mux    *http.ServeMux
	router *mux.Router
}

// NewServer returns a new Sever that will listen on addr.
func NewServer(addr string) Server {
	var s Server
	s.Addr = addr
	s.mux = http.NewServeMux()
	s.router = mux.NewRouter()
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

// A Route represents a REST route.
type Route mux.Route

// Route registers a new route for the specified URL path and allows to
// register the handlers pipe that will be used to handle the requests
// to that resource. For more information about the path syntax see the
// documentation of the package "github.com/gorilla/mux".
//
// The handlers pipe is executed sequentially until a HTTP header is
// explicitally written in the response. Besides that, mandatory handlers
// are always executed.
func (s Server) Route(path string, handlers ...http.Handler) *Route {
	return (*Route)(s.router.Handle(path, newPipe(handlers...)))
}

// RouteDefault sets the default route. This handler is used if any other
// route matches the request URI.
func (s Server) RouteDefault(handlers ...http.Handler) {
	s.router.NotFoundHandler = newPipe(handlers...)
}

// Methods allows to filter which HTTP methods will be handled by a given
// route. e.g.: "GET", "POST", "PUT", etc.
func (r *Route) Methods(methods ...string) *Route {
	mr := (*mux.Route)(r)
	return (*Route)(mr.Methods(methods...))
}

// Vars returns the route variables for the current request, if any.
func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}
