# httpw

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/krostar/httpw)
[![Licence](https://img.shields.io/github/license/krostar/httpw.svg?style=for-the-badge)](https://tldrlegal.com/license/mit-license)
![Latest version](https://img.shields.io/github/tag/krostar/httpw.svg?style=for-the-badge)

[![Build Status](https://img.shields.io/travis/krostar/httpw/master.svg?style=for-the-badge)](https://travis-ci.org/krostar/httpw)
[![Code quality](https://img.shields.io/codacy/grade/84f74110bd71455ea3a20b4163be7010/master.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/httpw/dashboard)
[![Code coverage](https://img.shields.io/codacy/coverage/84f74110bd71455ea3a20b4163be7010.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/httpw/dashboard)

Idiomatic HTTP handler signature that requires no magic.

## Motivation

On one hand, I hate when people are creating libraries that require
to change the signature of a http handler (`http.HandlerFunc`)
without any usage-related reasons (hi _gin_).
There is nearly nothing impossible to do with the actual `net/http` handler signature. On another hand, I have few problems related to the actual `net/http` handler signature:

-   there is no return argument, so when I call `rw.WriteHeader` on any condition
    I **always** forget to call `return` after ; I find this so not go idiomatic
    to require multiple lines to handle simple things
-   there is no simple way to ensure the same response and error format,
    the same way of handling them, on all handlers when multiple people
    write multiple handlers
-   tests require `httptest.NewRecorder` because nothing is returned,
    IMO it's boring

So I started thinking on removing the possibility of using directly the 
`http.ResponseWriter` and returning an argument that will fill it instead.

I also wanted to still be fully `net/http` compliant as I really don't want to
create yet another framework that handle and force too much things like routing,
handlers, tests, ...

The solution I found and iterate on for multiple weeks is the following:
`func(r *http.Request) (*httpw.Response, error)`. With the following signature
any handlers can handle the request as they would normally do, yet does not handle
how response is written, nor how errors are handled. This is left to a callback
that write the response, and another callback to handle the error, if any (for
example adding fields in a logger, set a tracing span to any status, ...).

## Usage and example

All `httpw.Handler` handlers can be used as normal `net/http` handlers with
`httpw.Wrap` method.

```go
// doSomething will return 200 if everything succeeds
// or 500 if something went wrong
func doSomething(r *http.Request) (*httpw.R, error) {
    var userID = getUserID(r)

    if err := doSomething(r.Context(), userID); err != nil {
        return nil, errors.Wrapf(err, "unable to do something with user %s", userID)
    }

    return &httpw.R{Status: http.StatusOK}
}
```

```go
func doSomething(r *http.Request) (*httpw.R, error) {
    var userID = getUserID(r)

    user, err := getSomething(r.Context(), userID)
    if err != nil {
        var e httpw.E
        if nfErr, ok := err.(NotFoundError); ok && nfErr() {
            e.Status = http.Status
        }
        return nil, e
    }

    return &httpw.R{
        Status: http.StatusOK,
        Data:   user,
    }
}
```

By default, if a standard error is returned, the request status will be set to
`http.StatusInternalServerError`. This behaviour can be changed by passing options
to the wrapper. Here are all the differents options:

```go
func setupRoutes() {
    var wrapper = httpw.New(
        httpw.WithDefaultErrorStatus(http.StatusServiceUnavailable),
        httpw.WithOnErrorCallback(logOnError),
        httpw.WithDataMarshaler(json.Marshal),
    )

    router.Get("/something", wrapper.Wrap(doSomething))
}
```

One last things: there are two aliases, `httpw.R` for `httpw.Response` and
`httpw.E` for `httpw.Error`

## License

This project is under the MIT licence, please see the LICENCE file.
