// Copyright 2014 The orujo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package sessions implements the bult-in sessions handler of
orujo.
*/
package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jroimartin/orujo"
)

// A SessionHandler is a orujo built-in handler that provides
// session management features.
type SessionHandler struct {
	sessionName string
	cookieStore *sessions.CookieStore
}

// An Options object contains the properties of the cookie that
// will be used to store the user session. See the
// documentation of the package "github.com/gorilla/sessions"
// for details.
type Options sessions.Options

// NewSessionHandler returns a new SessionHandler. name is
// used to set the sesion name. secret is used to specify
// the key used to authenticate the session.
func NewSessionHandler(name string, secret []byte) SessionHandler {
	return SessionHandler{
		sessionName: name,
		cookieStore: sessions.NewCookieStore(secret),
	}
}

// SetOptions sets the options for the cookie that will store
// the user session.
func (h SessionHandler) SetOptions(opts *Options) {
	h.cookieStore.Options = (*sessions.Options)(opts)
}

// Options retrieves the options of the cookie that stores
// the user session.
func (h SessionHandler) Options() *Options {
	return (*Options)(h.cookieStore.Options)
}

// ServeHTTP generates a new session id if the user does not own one
// yet.
func (h SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := h.SessionID(r); err == nil {
		return
	}
	session, err := h.cookieStore.Get(r, h.sessionName)
	if err != nil {
		internalServerError(w)
		orujo.RegisterError(w, err)
	}
	sessionID, err := randomString()
	if err != nil {
		internalServerError(w)
		orujo.RegisterError(w, err)
	}
	session.Values["id"] = sessionID
	if err := session.Save(r, w); err != nil {
		internalServerError(w)
		orujo.RegisterError(w, err)
	}
}

// SessionID retrieves the session id of the current user.
func (h SessionHandler) SessionID(r *http.Request) (string, error) {
	cookie, err := h.cookieStore.Get(r, h.sessionName)
	if err != nil {
		return "", err
	}
	sessionID, ok := cookie.Values["id"].(string)
	if !ok {
		return "", errors.New("Session ID is not a string")
	}
	return sessionID, nil
}

func internalServerError(w http.ResponseWriter) {
	status := http.StatusInternalServerError
	http.Error(w, http.StatusText(status), status)
}

func randomString() (string, error) {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	str := base64.StdEncoding.EncodeToString(buf)
	return str, nil
}
