[![Builds](https://github.com/ssrathi/go-attr/workflows/Build/badge.svg?branch=master)](https://github.com/ssrathi/go-attr/actions?query=branch%3Amaster+workflow%3ABuild)
[![Go Report Card](https://goreportcard.com/badge/github.com/ssrathi/go-attr)](https://goreportcard.com/report/github.com/ssrathi/go-attr)
[![GoDoc](https://godoc.org/github.com/ssrathi/go-attr?status.svg)](https://pkg.go.dev/github.com/ssrathi/go-attr)

# go-attr
Golang library to act on structure fields at runtime. Similar to Python getattr(), setattr(), hasattr() APIs.

This package provides user friendly helper APIs built on top of the Golang "reflect" library. Reflect library is tricky to use due to its low level nature and results in a panic if an incorrect input is provided. This package provides high level abstractions on such tricky APIs in a user friendly manner.

## Installation
```
go get -u github.com/ssrathi/go-attr
```

Or install manually
```
git clone https://github.com/ssrathi/go-attr
cd go-attr
go install
```

## Usage
See full documentation at https://pkg.go.dev/github.com/ssrathi/go-attr.

```go
  import attr "github.com/ssrathi/go-attr"

  type User struct {
    Username string `json:"username" db:"uname"`
    Age      int `json:"age" meta:"important"`
    password string
  }
  user := User{"srathi", 30, "my_secret_123"}

  // NOTE:  Handle error if present in the examples below.

  // HasField(): Check if a field name is part of a struct object.
  ok, err := attr.HasField(&user, "FirstName")
  fmt.Printf("FirstName found: %v\n", ok)

  // SetField(): Set a new value to an existing field of a struct object.
  err = attr.SetField(&user, "Username", "new-username")
  fmt.Printf("New username: %s\n", user.Username)

  // GetField(): Get the current value of a struct object.
  val, err = attr.GetField(&user, "Username")
  fmt.Printf("Username: %s\n", user.Username)

  // FieldValues(): Get the values of all the fields.
  fieldValues, err := attr.FieldValues(&user)
  for name, val := range fieldValues {
    fmt.Printf("%s: %v\n", name, val)
  }

  // GetFieldTag(): Get the value of tag "meta" of field "Age".
  tagValue, err := attr.GetFieldTag(&user, "Age", "meta")
  fmt.Printf("'meta' tag value of 'Age': %s\n", tagValue)

  // TagValues(): Get the value of tag 'json' from all public fields of a struct.
  // Tag value is blank ("") if it is not part of a public field.
  tagVals, err := attr.TagValues(&user, "json")
  for fieldName, tagVal := range tagVals {
    fmt.Printf("%s: %v\n", fieldName, tagVal)
  }
```

## Contributing

Contributions are most welcome! Please follow the steps below to send
pull requests with your changes.

* Fork this repository and create a feature branch in it.
* Push a commit with your changes.
* Create a new pull request.
* Create a new issue and link the pull request to it.
