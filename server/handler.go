// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import "net/http"

type Handler interface {
	SetMandatory(v bool)
	Mandatory() bool
	ServeHTTP(rw http.ResponseWriter, r *http.Request)
}

func M(h Handler) Handler {
	h.SetMandatory(true)
	return h
}

type handlerWrapper struct {
	hh http.Handler
	mandatory bool
}

func (hw *handlerWrapper) Mandatory() bool {
	return hw.mandatory
}

func (hw *handlerWrapper) SetMandatory(v bool) {
	hw.mandatory = v
}

func (hw *handlerWrapper) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	hw.hh.ServeHTTP(rw, r)
}

func H(hh http.Handler) Handler {
	return &handlerWrapper{hh: hh, mandatory: false}
}
