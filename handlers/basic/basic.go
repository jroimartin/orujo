// Copyright 2014 The gorest Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package basic implements basic auth mechanisms for gorest.
*/
package basic

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/jroimartin/gorest"
)

// A BasicHandler is a gorest built-in handler that provides
// basic authentication.
type BasicHandler struct {
	realm    string
	username string
	password string
}

// NewBasicHandler returns a new BasicHandler.
func NewBasicHandler(realm, username, password string) BasicHandler {
	return BasicHandler{
		realm:    realm,
		username: username,
		password: password,
	}
}

// ServeHTTP validates the username and password provided by the user.
func (h BasicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	isValid, provUser := h.isValidAuth(r.Header.Get("Authorization"))
	if isValid {
		return
	}
	w.Header().Set("WWW-Authenticate", "Basic realm=\""+h.realm+"\"")
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintln(w, http.StatusText(http.StatusUnauthorized))
	errorStr := fmt.Errorf("Invalid username or password (username: %s)", provUser)
	gorest.RegisterError(w, errorStr)
}

func (h BasicHandler) isValidAuth(auth string) (valid bool, username string) {
	b64auth := strings.Split(auth, " ")
	if len(b64auth) != 2 || b64auth[0] != "Basic" {
		return false, "unknown"
	}
	creds, err := base64.StdEncoding.DecodeString(b64auth[1])
	if err != nil {
		return false, "unknown"
	}
	userpass := strings.Split(string(creds), ":")
	if len(userpass) != 2 {
		return false, "unknown"
	}

	provUserSha256 := sha256.Sum256([]byte(userpass[0]))
	userSha256 := sha256.Sum256([]byte(h.username))
	validUser := subtle.ConstantTimeCompare(provUserSha256[:], userSha256[:]) == 1

	provPassSha256 := sha256.Sum256([]byte(userpass[1]))
	passSha256 := sha256.Sum256([]byte(h.password))
	validPass := subtle.ConstantTimeCompare(provPassSha256[:], passSha256[:]) == 1

	return validUser && validPass, userpass[0]
}
