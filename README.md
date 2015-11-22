# Orujo [![GoDoc](https://godoc.org/github.com/jroimartin/orujo?status.svg)](https://godoc.org/github.com/jroimartin/orujo)

## Introduction

Orujo solves a common problem, which is the execution of several middlewares per
route. It has been designed to work seamlessly with the standard net/http
library. 

## Routes and Pipes

So, how can I link several actions/middlewares with a route? the answer is
"pipes". Let me show this with a single example:

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/jroimartin/orujo"
)

func main() {
    http.Handle("/hello", orujo.NewPipe(
        http.HandlerFunc(helloHandler),
        orujo.M(http.HandlerFunc(worldHandler)),
    ))
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprint(w, "Hello, ")
}

func worldHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "world")
}
```

In this example, we are linking the route `/hello/` with the following pipe
of handlers:

```
helloHandler --> M(worldHandler)
```

A pipe is a sequence of handlers that will be executed until one of the handlers
explicitly calls the function ResponseWriter.WriteHeader(). From that moment
only mandatory handlers get executed. In this example, the only mandatory
handler in the pipe would be "worldHandler", which was marked as mandatory
using the helper function M().

## http.Handler

One of the main goals behind Orujo is standarization. Due to this, the handlers
accepted by Orujo must satisfy the interface http.Handler and, of course, the
returned pipe also implements the interface http.Handler. This way, everything
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
