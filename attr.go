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
// Example
//
// A quick example to see it in action.
// Full documentation is at https://pkg.go.dev/github.com/ssrathi/go-attr.
//
// 	import attr "github.com/ssrathi/go-attr"
//
// 	user := User{
// 		Username:  "srathi",
// 		FirstName: "Shyamsunder",
// 	}
//
// 	err = attr.SetValue(&user, "Username", "srathi-alt")
// 	fmt.Printf("New username: %s\n", user.Username) // prints "srathi-alt"
package attr

import (
	"errors"
	"reflect"
)

// Error values.
var (
	ErrNoField         = errors.New("Specified field is not present in the struct")
	ErrNotPtr          = errors.New("Specified struct is not passed by pointer")
	ErrNotStruct       = errors.New("Given object is not a struct or a pointer to a struct")
	ErrUnexportedField = errors.New("Specified field is not an exported or public field")
	ErrMismatchValue   = errors.New("Specified value to set is of a different type")
)

// GetValue returns the value of a given field of a structure given by 'obj'.
// 'obj' can be passed by value or by pointer.
// Only exported (public) field values can be found (else ErrUnexportedField is raised).
//
// If the field is not found, then an error is returned.
func GetValue(obj interface{}, fieldName string) (interface{}, error) {
	objValue, err := getReflectValue(obj)
	if err != nil {
		return nil, err
	}

	fieldValue := objValue.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return nil, ErrNoField
	}

	if !fieldValue.CanInterface() {
		return nil, ErrUnexportedField
	}

	return fieldValue.Interface(), nil
}

// Has returns a boolean indicating if the given field name is found in
// the given struct obj.
func Has(obj interface{}, fieldName string) (bool, error) {
	objValue, err := getReflectValue(obj)
	if err != nil {
		return false, err
	}

	structType := objValue.Type()
	_, found := structType.FieldByName(fieldName)
	return found, nil
}

// SetValue sets the given value to the fieldName field in the given struct 'obj'.
// Only exported (public) fields can be set using this API.
//
// NOTE: 'obj' struct must be passed by pointer for this API to work. Passing by
// value results in ErrPassedByValue.
func SetValue(obj interface{}, fieldName string, newValue interface{}) error {
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

// Names returns a slice of all field names of a given struct.
// Only the exportable (public) field names are returned.
func Names(obj interface{}) ([]string, error) {
	objValue, err := getReflectValue(obj)
	if err != nil {
		return nil, err
	}

	fieldNames := []string{}
	objType := objValue.Type()
	for i := 0; i < objValue.NumField(); i++ {
		fieldType := objType.Field(i)
		fieldValue := objValue.Field(i)

		if fieldValue.CanInterface() {
			fieldNames = append(fieldNames, fieldType.Name)
		}
	}

	return fieldNames, nil
}

// Values returns a map of all field names with the value of each field.
// Only the exportable (public) field name-value pairs are returned.
func Values(obj interface{}) (map[string]interface{}, error) {
	objValue, err := getReflectValue(obj)
	if err != nil {
		return nil, err
	}

	valueMap := map[string]interface{}{}
	objType := objValue.Type()
	for i := 0; i < objValue.NumField(); i++ {
		fieldType := objType.Field(i)
		fieldValue := objValue.Field(i)

		if fieldValue.CanInterface() {
			valueMap[fieldType.Name] = fieldValue.Interface()
		}
	}

	return valueMap, nil
}

// GetTag returns the value of a specified tag on a specified struct field.
// Specified field must be an exportable (public) filed of the struct.
func GetTag(obj interface{}, fieldName, tagKey string) (string, error) {
	objValue, err := getReflectValue(obj)
	if err != nil {
		return "", err
	}

	structType := objValue.Type()
	field, found := structType.FieldByName(fieldName)
	if !found {
		return "", ErrNoField
	}

	if field.PkgPath != "" {
		return "", ErrUnexportedField
	}

	return field.Tag.Get(tagKey), nil
}

// Tags returns a map of all the tag values of a given tag key from all
// the exported (public) struct fields.
func Tags(obj interface{}, tagKey string) (map[string]string, error) {
	objValue, err := getReflectValue(obj)
	if err != nil {
		return nil, err
	}

	tagMap := map[string]string{}
	objType := objValue.Type()
	for i := 0; i < objValue.NumField(); i++ {
		fieldType := objType.Field(i)
		fieldValue := objValue.Field(i)

		if fieldValue.CanInterface() {
			tagMap[fieldType.Name] = fieldType.Tag.Get(tagKey)
		}
	}

	return tagMap, nil
}

// GetKind returns the "kind" of a specified public struct field. "Kind" is
// the in-built type of a variable, such as Uint64, Slice, Struct, Ptr, etc.
func GetKind(obj interface{}, fieldName string) (string, error) {
	objValue, err := getReflectValue(obj)
	if err != nil {
		return "", err
	}

	fieldValue := objValue.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return "", ErrNoField
	}

	if !fieldValue.CanInterface() {
		return "", ErrUnexportedField
	}

	return fieldValue.Kind().String(), nil
}

// Kinds returns the 'kind' of all the public fields of a struct. "Kind" is
// the in-built type of a variable, such as Uint64, Slice, Struct, Ptr, etc.
func Kinds(obj interface{}) (map[string]string, error) {
	objValue, err := getReflectValue(obj)
	if err != nil {
		return nil, err
	}

	kindMap := map[string]string{}
	objType := objValue.Type()
	for i := 0; i < objValue.NumField(); i++ {
		fieldType := objType.Field(i)
		fieldValue := objValue.Field(i)

		if fieldValue.CanInterface() {
			kindMap[fieldType.Name] = fieldValue.Kind().String()
		}
	}

	return kindMap, nil
}

// getReflectValue gets a reflect-value of a given struct. If it is a pointer
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
