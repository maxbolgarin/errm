# errm

[![GoDoc][doc-img]][doc] [![Build][ci-img]][ci] [![GoReport][report-img]][report]

Package `errm` is wrapper on [eris](https://github.com/rotisserie/eris) for convinient usage of errors with structrual fields and stack trace

Install: `go get github.com/maxbolgarin/errm`

- [Why you should try](#why-you-should-try)
- [How to use](#how-to-use)
- [Contributing](#contributing)

## Why you should try

There are familiar methods like `New`, `Errorf`, `Wrap` and others that works as expected. But there are two breaking futures:

1. There is a stack trace in every error, thanks for the `eris`
2. You can add `field=value` pairs to make an error message more convinient to handle in future

Which option is better for further search and analysis?

* err1: `cannot start server 'orders' at address: :7000: port is in use`
* err2: `cannot start server=orders address=:7000: port is in use`

The second one can be easily parsed and it is more friedly to read. Here is the code for these two examples:

```go
err1 := fmt.Errorf("cannot start server '%s' at address: %s: %w", name, addr, err)
err2 := errm.Wrap(err, "cannot start", "server", name, "address", addr)
```

## How to use


### New error

```go
err := errm.New("some-err", "field", "value", "field2", []any{123, 321}, "field3", 123, "field4")
fmt.Println(err) 

// some-err field=value field2=[123 321] field3=123
```

### Wrap error

```go
notFoundErr := errm.New("not found")

name := "database"
err2 := errm.Wrapf(err1, "another error with %s", name, "address", "127.0.0.1")

if errm.Is(err2, notFoundErr) {
    fmt.Println(err2)
}

// another error with database address=127.0.0.1: not found
```

### Error List

```go
randomErr := errm.Errorf("some random error with %s", "unwanted behaviour")
notFoundErr := errm.New("not found")

errList := errm.NewList()
errList.Add(randomErr)
errList.New("some error", "retry", 1)
errList.Wrap(notFoundErr, "database error")

if errList.Has(notFoundErr) {
    finalErr := errList.Err()
    fmt.Println(errm.Wrap(finalErr.Error(), "multi error")) 
}

// multi error: some random error with unwanted behaviour; some error retry=1; database error: not found
```

### Error Set

```go
errSet := errm.NewSet()

// Add the same error for three times
notFoundErr := errm.New("not found")
errSet.Add(notFoundErr)
errSet.Add(notFoundErr)
errSet.Add(notFoundErr)
fmt.Println(errSet.Len()) // 1

// Add two errors with the same message
errSet.Add(errm.New("A"))
errSet.Add(errm.New("A"))
fmt.Println(errSet.Len()) // 2

// Add an another error but with the same message
err := errors.New("A")
errSet.Add(err)
fmt.Println(errSet.Len()) // 2

// Add wrapped -> got another error message -> another error
errSet.Add(errm.Wrap(notFoundErr, "another error"))
fmt.Println(errSet.Len()) // 3

```

## Contributing

If you'd like to contribute to `errm`, make a fork and submit a pull request!

Released under the [MIT License]

[MIT License]: LICENSE.txt
[doc-img]: https://pkg.go.dev/badge/github.com/maxbolgarin/errm
[doc]: https://pkg.go.dev/github.com/maxbolgarin/errm
[ci-img]: https://github.com/maxbolgarin/errm/actions/workflows/go.yaml/badge.svg
[ci]: https://github.com/maxbolgarin/errm/actions
[report-img]: https://goreportcard.com/badge/github.com/maxbolgarin/errm
[report]: https://goreportcard.com/report/github.com/maxbolgarin/errm
