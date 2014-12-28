# Orujo [![GoDoc](https://godoc.org/github.com/jroimartin/orujo?status.svg)](https://godoc.org/github.com/jroimartin/orujo)

## Introduction

Orujo is a minimalist web framework written in Go, which has been designed
to work seamlessly with the standard net/http library. 

## What does an Orujo application look like? 

The following snippet shows a very simple "hello world" program written using
Orujo. It registers the route "/" and links it with a http.Handler that will
return the string "Hello world!" when it is requested by the user.

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jroimartin/orujo"
)

func main() {
	s := orujo.NewServer("localhost:8080")

	s.Route(`^/$`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world!")
	}))

	log.Fatalln(s.ListenAndServe())
}
```

## Routes and Pipes

Some of the first questions that maybe come to your mind are:

* How can I define a route?
* How can I link several actions/middlewares with a route?

Let me answer both questions with a single example:

```go
s := orujo.NewServer("localhost:8080")
s.Route(`^/private/.*`,
	authHandler,
	myHandler,
	orujo.M(logHandler),
).Methods("GET", "POST")
```

In this example we are registering a new route, defined by the
[regular expression](http://golang.org/pkg/regexp/) `^/private/.*` and the valid
HTTP methods "GET" and "POST". Besides that, this route is also linked to the
following pipe of handlers:

```
authHandler --> myHandler --> M(logHandler)
```

A pipe is a sequence of handlers that will be executed until one of the handlers
explicitly calls the function ResponseWriter.WriteHeader(). From that moment
only mandatory handlers get executed. In this example, the only mandatory
handler in the pipe would be "logHandler", which was marked as mandatory using
the helper function M().

## http.Handler

One of the main goals behind Orujo is standarization. Due to this, the handlers
accepted by Orujo must satisfy the interface http.Handler. This way, everything
that already works with the Go's standard library must work with Orujo.

```go
func (h LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	...
}
```

A repository with handlers ready to use with Orujo can be found
[here](https://github.com/jroimartin/orujo-handlers).

```go
import "github.com/jroimartin/orujo-handlers/<name>"
```

## Installation

`go get github.com/jroimartin/orujo`

## More information

`godoc github.com/jroimartin/orujo`
