// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jroimartin/gorest"
)

type SessionHandler struct {
	sessionName string
	cookieStore *sessions.CookieStore
	mandatory   bool
}

type Options sessions.Options

func NewSessionHandler(name string, secret []byte) *SessionHandler {
	h := &SessionHandler{}
	h.sessionName = name
	h.cookieStore = sessions.NewCookieStore(secret)
	return h
}

func (h *SessionHandler) SetOptions(opts *Options) {
	h.cookieStore.Options = (*sessions.Options)(opts)
}

func (h *SessionHandler) Options() *Options {
	return (*Options)(h.cookieStore.Options)
}

func (h *SessionHandler) SetMandatory(v bool) {
	h.mandatory = v
}

func (h *SessionHandler) Mandatory() bool {
	return h.mandatory
}

func (h *SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := h.SessionId(r); err == nil {
		return
	}
	session, err := h.cookieStore.Get(r, h.sessionName)
	if err != nil {
		internalServerError(w)
		gorest.RegisterError(w, err)
	}
	sessionId, err := randomString()
	if err != nil {
		internalServerError(w)
		gorest.RegisterError(w, err)
	}
	session.Values["id"] = sessionId
	if err := session.Save(r, w); err != nil {
		internalServerError(w)
		gorest.RegisterError(w, err)
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
