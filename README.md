[![Builds](https://github.com/ssrathi/go-attr/workflows/Build/badge.svg?branch=master)](https://github.com/ssrathi/go-attr/actions?query=branch%3Amaster+workflow%3ABuild)
[![Go Report Card](https://goreportcard.com/badge/github.com/ssrathi/go-attr)](https://goreportcard.com/report/github.com/ssrathi/go-attr)
[![GoDoc](https://godoc.org/github.com/ssrathi/go-attr?status.svg)](https://godoc.org/github.com/ssrathi/go-attr)

# go-attr
Golang library to act on structure fields at runtime. Similar to Python getattr(), setattr(), hasattr() APIs.

This package provides user friendly helper APIs built on top of the Golang "reflect" library. Reflect library is tricky to use due to its low level nature and results in a panic if an incorrect input is provided. This package provides high level abstractions on such tricky APIs in a user friendly manner.

## Installation
```
go get github.com/ssrathi/go-attr
```

## Usage
See full documentation at https://godoc.org/github.com/ssrathi/go-attr
```go
  import (
    attr "github.com/ssrathi/go-attr"
  )

  type User struct {
    Username  string
    FirstName string
  }

  user := User{
    Username:  "srathi",
    FirstName: "Shyamsunder",
  }

  // Check if a field name is part of a struct object.
  ok, err := attr.HasField(&user, "FirstName")
  if err != nil {
    // Handle error.
  }
  fmt.Printf("FirstName found: %v\n", ok)

  // Set a new value to an existing field of a struct object.
  err = attr.SetField(&user, "Username", "new-username")
  if err != nil {
    // Handle error.
  }
  fmt.Printf("New username: %s\n", user.Username)

  // Get the current value of a struct object.
  val, err = attr.GetField(&user, "Username")
  if err != nil {
    // Handle error.
  }
  fmt.Printf("Username: %s\n", user.Username)

  // Get the values of all the fields.
  fieldValues, err := attr.FieldValues(&user)
  if err != nil {
    // Handle error.
  }
  for name, val := range fieldValues {
    fmt.Printf("%s: %v\n", name, val)
  }
```

## Contributing

Contributions are most welcome! Please follow the steps below to send
pull requests with your changes.

* Fork this repository and create a feature branch in it.
* Push a commit with your changes.
* Create a new pull request.
* Create a new issue and link the pull request to it.
