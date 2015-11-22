// Copyright 2014 The orujo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package basic implements basic auth mechanisms for orujo.
*/
package basic

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/jroimartin/orujo"
)

// A BasicHandler is a orujo built-in handler that provides
// basic authentication.
type BasicHandler struct {
	realm    string
	username string
	password string

	// ErrorMsg can be used to set a custom error message. The parameter
	// provuser contains the username provided by the user.
	ErrorMsg func(w http.ResponseWriter, provuser string)
}

// NewBasicHandler returns a new BasicHandler.
func NewBasicHandler(realm, username, password string) BasicHandler {
	return BasicHandler{
		realm:    realm,
		username: username,
		password: password,
		ErrorMsg: defaultErrorMsg,
	}
}

// ServeHTTP validates the username and password provided by the user.
func (h BasicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isValid, provUser := h.isValidAuth(r)
	if isValid {
		return
	}
	w.Header().Set("WWW-Authenticate", "Basic realm=\""+h.realm+"\"")
	w.WriteHeader(http.StatusUnauthorized)
	h.ErrorMsg(w, provUser)
	errorStr := fmt.Errorf("Invalid username or password (username: %s)", provUser)
	orujo.RegisterError(w, errorStr)
}

// defaultErrorMsg writes the default unauthorized response to the
// http.ResponseWriter
func defaultErrorMsg(w http.ResponseWriter, provuser string) {
	fmt.Fprintln(w, http.StatusText(http.StatusUnauthorized))
}

func (h BasicHandler) isValidAuth(r *http.Request) (valid bool, username string) {
	provUser, provPass, ok := r.BasicAuth()
	if !ok {
		return false, "unknown"
	}

	provUserSha256 := sha256.Sum256([]byte(provUser))
	userSha256 := sha256.Sum256([]byte(h.username))
	validUser := subtle.ConstantTimeCompare(provUserSha256[:], userSha256[:]) == 1

	provPassSha256 := sha256.Sum256([]byte(provPass))
	passSha256 := sha256.Sum256([]byte(h.password))
	validPass := subtle.ConstantTimeCompare(provPassSha256[:], passSha256[:]) == 1

	return validUser && validPass, provUser
}
