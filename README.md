[![Builds](https://github.com/ssrathi/go-attr/workflows/Build/badge.svg?branch=master)](https://github.com/ssrathi/go-attr/actions?query=branch%3Amaster+workflow%3ABuild)
[![Go Report Card](https://goreportcard.com/badge/github.com/ssrathi/go-attr)](https://goreportcard.com/report/github.com/ssrathi/go-attr)
[![codecov.io](https://codecov.io/github/ssrathi/go-attr/coverage.svg?branch=master)](https://codecov.io/github/ssrathi/go-attr?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/ssrathi/go-attr.svg)](https://pkg.go.dev/github.com/ssrathi/go-attr)

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
    Age      int    `json:"age" meta:"important"`
    password string
  }
  
  user := User{"srathi", 30, "my_secret_123"}
  // NOTE:  Handle error if present in the examples below.
  // Also, all the APIs work on the exported (public) fields of a struct.
```
### SetValue()

**Set a new value to an existing field of a struct object.**
```go
  // Struct must be passed by pointer to set its field.
  err = attr.SetValue(&user, "Username", "new-username")
  fmt.Printf("New username: %s\n", user.Username)
```
### GetValue()

**Get the current value of a struct object.**
```go
  val, err = attr.GetValue(&user, "Username")
  fmt.Printf("Username: %s\n", user.Username)
```
### Has()

**Check if a field name is part of a struct object.**
```go
  ok, err := attr.Has(&user, "FirstName")
  fmt.Printf("FirstName found: %v\n", ok)
```
### Names()

**Get the names of all the struct fields.**
```go
  fieldNames, err := attr.Names(&user)
  fmt.Printf("field names: %v\n", fieldNames)
```
### Values()

**Get the values of all the struct fields.**
```go
  fieldValues, err := attr.Values(&user)
  for name, val := range fieldValues {
    fmt.Printf("%s: %v\n", name, val)
  }
```
### GetTag()

**Get the value of a specific tag of a specific field in a struct.**
```go
  tagValue, err := attr.GetTag(&user, "Age", "meta")
  fmt.Printf("'meta' tag value of 'Age': %s\n", tagValue)
```
### Tags()

**Get the value of specific tag from all the public fields of a struct.**
```go
  // Tag value is blank ("") if it is not part of a public field.
  tagVals, err := attr.Tags(&user, "json")
  for fieldName, tagVal := range tagVals {
    fmt.Printf("%s: %v\n", fieldName, tagVal)
  }
```
### GetKind()

**Get the "kind" (type) of a specified struct field.**
```go
  // For example, "var Age int" is of kind 'int'.
  kind, err := attr.GetKind(&user, "Age")
  fmt.Printf("Kind of 'Age': %s\n", kind)
```
### Kinds()

**Get the "kind" (type) of all the struct fields.**
```go
  // For example, "var Age int" is of kind 'int'.
  kinds, err := attr.Kinds(&user)
  for name, kind := range kinds {
    fmt.Printf("%s: %s\n", name, kind)
  }
```

## Contributing

Contributions are most welcome! Please follow the steps below to send
pull requests with your changes.

* Fork this repository and create a feature branch in it.
* Push a commit with your changes.
* Create a new pull request.
* Create a new issue and link the pull request to it.
