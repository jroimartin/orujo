// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gorest

import "net/http"

// Objects imeplementing the Handler interface can be used
// as part of a handlers pipe.
type Handler interface {
	// SetMandatory allows to set the "mandatory" property of
	// a Handler.
	SetMandatory(v bool)

	// Mandatory returns if a Handler is mandatory or not.
	Mandatory() bool

	// ServeHTTP implements the Handler's functionality.
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// M is a helper to easily set a Handler as "mandatory".
func M(h Handler) Handler {
	h.SetMandatory(true)
	return h
}

type handlerWrapper struct {
	hh        http.Handler
	mandatory bool
}

func (hw *handlerWrapper) Mandatory() bool {
	return hw.mandatory
}

func (hw *handlerWrapper) SetMandatory(v bool) {
	hw.mandatory = v
}

func (hw *handlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hw.hh.ServeHTTP(w, r)
}

// H is a helper to convert a http.Handler to a REST Handler.
func H(hh http.Handler) Handler {
	return &handlerWrapper{hh: hh, mandatory: false}
}
