// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

type SessionHandler struct {
	sessionName string
	cookieStore *sessions.CookieStore
	mandatory   bool
}

func NewSessionHandler(name string, secret []byte) *SessionHandler {
	h := &SessionHandler{}
	h.sessionName = name
	h.cookieStore = sessions.NewCookieStore(secret)
	return h
}

func (h *SessionHandler) SetMandatory(v bool) {
	h.mandatory = v
}

func (h *SessionHandler) Mandatory() bool {
	return h.mandatory
}

func (h *SessionHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	if _, err := h.SessionId(r); err == nil {
		return
	}
	session, err := h.cookieStore.Get(r, h.sessionName)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	sessionId, err := randomString()
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	session.Values["id"] = sessionId
	if err := session.Save(r, rw); err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *SessionHandler) SessionId(r *http.Request) (string, error) {
	cookie, err := h.cookieStore.Get(r, h.sessionName)
	if err != nil {
		return "", err
	}
	sessionId, ok := cookie.Values["id"].(string)
	if !ok {
		return "", errors.New("Session ID is not a string")
	}
	return sessionId, nil
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
