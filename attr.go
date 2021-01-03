/*
 * Author: Shyamsunder Rathi (shyam29@gmail.com)
 *
 * License: MIT (See License file for full text).
 */

// Package attr implements Python style APIs to access and modify structure
// fields and tags at the runtime, similar to Python APIs getattr(), setattr(),
// hasattr(), etc.
//
// This package provides user friendly helper APIs built on top of the Golang
// "reflect" library. Reflect library is tricky to use due to its low level nature
// and results in a panic if an incorrect input is provided. This package provides
// high level abstractions on such tricky APIs in a user friendly manner.
//
// 	import attr "github.com/ssrathi/go-attr"
//
// 	user := User{
// 		Username:  "srathi",
// 		FirstName: "Shyamsunder",
// 	}
//
// 	ok, err := attr.HasField(&user, "FirstName")
// 	fmt.Printf("FirstName found: %v\n", ok)
//
// 	err = attr.SetField(&user, "Username", "new-username")
// 	fmt.Printf("New username: %s\n", user.Username)
//
// 	val, err = attr.GetField(&user, "Username")
// 	fmt.Printf("Username: %s\n", user.Username)
package attr

import (
	"errors"
	"reflect"
)

var (
	// ErrNoField is raised if a non-existent struct field is passed.
	ErrNoField = errors.New("Specified field is not present in the struct")
	// ErrNotPtr is raised if struct is passed by value while trying to set a field.
	ErrNotPtr = errors.New("Specified struct is not passed by pointer")
	// ErrNotStruct is raised if a struct or a pointer to a struct is not used.
	ErrNotStruct = errors.New("Given object is not a struct or a pointer to a struct")
	// ErrUnexportedField is raised if 'set' is tried on a private field of a struct.
	ErrUnexportedField = errors.New("Specified field is not an exported or public field")
	// ErrMismatchValue is raised if 'set' is tried with a different type of value.
	ErrMismatchValue = errors.New("Specified value to set is of a different type")
)

// GetField returns the value of a given field of a structure given by 'obj'.
// 'obj' can be passed by value or by pointer.
// Only exported (public) field values can be found (else ErrUnexportedField is raised).
//
// If the field is not found, then an error is returned.
func GetField(obj interface{}, fieldName string) (interface{}, error) {
	objValue, err := getReflectValue(obj)
	if err != nil {
		return nil, err
	}

	structType := objValue.Type()
	field, found := structType.FieldByName(fieldName)
	if !found {
		return nil, ErrNoField
	}

	if !isExportedField(field) {
		return nil, ErrUnexportedField
	}

	fieldValue := objValue.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return nil, ErrNoField
	}

	return fieldValue.Interface(), nil
}

// HasField returns a boolean indicating if the given field name is found in
// the given struct obj.
func HasField(obj interface{}, fieldName string) (bool, error) {
	objValue, err := getReflectValue(obj)
	if err != nil {
		return false, err
	}

	structType := objValue.Type()
	_, found := structType.FieldByName(fieldName)
	return found, nil
}

// SetField sets the given value to the fieldName field in the given struct 'obj'.
// Only exported (public) fields can be set using this API.
//
// NOTE: 'obj' struct must be passed by pointer for this API to work. Passing by
// value results in ErrPassedByValue.
func SetField(obj interface{}, fieldName string, newValue interface{}) error {
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() != reflect.Ptr {
		return ErrNotPtr
	}

	objValue = objValue.Elem()
	if objValue.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	fieldValue := objValue.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return ErrNoField
	}

	if fieldValue.Type() != reflect.TypeOf(newValue) {
		return ErrMismatchValue
	}

	if !fieldValue.CanSet() {
		return ErrUnexportedField
	}

	fieldValue.Set(reflect.ValueOf(newValue))
	return nil
}

// getReflectValue gets a reflect-value of a given structure. If it is a pointer
// to a struct, then it gives the reflect-value of the underlying structure.
//
// Returns an error if the given obj is not a struct or a pointer to a struct.
func getReflectValue(obj interface{}) (reflect.Value, error) {
	value := reflect.ValueOf(obj)

	if value.Kind() == reflect.Struct {
		return value, nil
	}

	if value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct {
		return value.Elem(), nil
	}

	var retval reflect.Value
	return retval, ErrNotStruct
}

// isExportedField returns if a structure field is Exported or Public.
// An exported or public field is one which begins with a "unicode capital" letter.
//
// "reflect" package cannot get values of unexported (private) fields of a struct.
func isExportedField(field reflect.StructField) bool {
	// PkgPath is the package path that qualifies a lower case (unexported)
	// field name. It is empty for upper case (exported) field names.
	// See https://golang.org/ref/spec#Uniqueness_of_identifiers
	return field.PkgPath == ""
}
