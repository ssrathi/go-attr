package attr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type User struct {
	Username string
	password string
	Age      int
}

var user User = User{"srathi", "my_secret_123", 30}

func TestGetField(t *testing.T) {
	want := user.Username
	got, err := GetField(user, "Username")
	require.Nil(t, err)
	require.Equal(t, got, want, "Username mismatch")

	wantErr := ErrUnexportedField
	_, gotErr := GetField(user, "password")
	require.Equal(t, gotErr, wantErr, "Able to get an unexported field value")
}

func ExampleGetField() {
	testUser := User{
		Username: "srathi",
		Age:      30,
	}

	value, err := GetField(testUser, "Age")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("Value of Age: %v", value)
	// Output: Value of Age: 30
}

func TestHasField(t *testing.T) {
	want := true
	got, err := HasField(&user, "Age")
	require.Nil(t, err)
	require.Equal(t, got, want, "Age not found")

	want = false
	got, err = HasField(&user, "NonExistant")
	require.Nil(t, err)
	require.Equal(t, got, want, "NonExistant field found")
}

func ExampleHasField() {
	testUser := User{
		Username: "srathi",
		Age:      30,
	}

	ok, err := HasField(&testUser, "Age")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("Age found: %v\n", ok)

	ok, err = HasField(&testUser, "ABC")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("ABC found: %v\n", ok)

	// Output:
	// Age found: true
	// ABC found: false
}

func TestSetField(t *testing.T) {
	wantErr := ErrNotPtr
	gotErr := SetField(user, "Username", "new-srathi")
	require.Equal(t, gotErr, wantErr, "Able to set fields on a struct by value")

	wantErr = ErrUnexportedField
	gotErr = SetField(&user, "password", "new-password")
	require.Equal(t, gotErr, wantErr, "Able to set a private field of a struct")

	wantErr = ErrMismatchValue
	gotErr = SetField(&user, "Age", 40.5)
	require.Equal(t, gotErr, wantErr, "Able to set float value to an int field")

	oldVal := user.Age
	err := SetField(&user, "Age", 40)
	require.Nil(t, err)
	require.Equal(t, user.Age, 40, "New value not set in a struct")

	err = SetField(&user, "Age", oldVal)
	require.Nil(t, err)
	require.Equal(t, user.Age, oldVal, "New value not set in a struct")
}

func ExampleSetField() {
	testUser := User{
		Username: "srathi",
		password: "secret",
	}

	err := SetField(&testUser, "password", "new-secret")
	fmt.Printf("Error while setting a private field: %v\n", err)

	err = SetField(testUser, "Username", "new-username")
	fmt.Printf("Error while passing struct by value: %v\n", err)

	err = SetField(&testUser, "Username", 100)
	fmt.Printf("Error while setting 100 in username: %v\n", err)

	err = SetField(&testUser, "Username", "new-username")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("New username: %s\n", testUser.Username)

	// Output:
	// Error while setting a private field: Specified field is not an exported or public field
	// Error while passing struct by value: Specified struct is not passed by pointer
	// Error while setting 100 in username: Specified value to set is of a different type
	// New username: new-username
}

func TestFieldNames(t *testing.T) {
	// Only public fields are returned.
	want := []string{"Username", "Age"}
	got, err := FieldNames(&user)
	require.Nil(t, err)
	require.Equal(t, got, want, "Struct field list is not correct")
}

func ExampleFieldNames() {
	testUser := User{
		Username: "srathi",
		password: "secret",
		Age:      30,
	}

	fields, err := FieldNames(&testUser)
	if err != nil {
		// Handle error.
	}
	fmt.Printf("Fields: %v", fields)
	// Output: Fields: [Username Age]
}

func TestFieldValues(t *testing.T) {
	// Only the value of the public fields are returned.
	want := map[string]interface{}{"Username": "srathi", "Age": 30}
	got, err := FieldValues(&user)
	require.Nil(t, err)
	require.Equal(t, got, want, "Struct field values are not correct")
}

func ExampleFieldValues() {
	testUser := User{
		Username: "srathi",
		password: "secret",
		Age:      30,
	}

	values, err := FieldValues(&testUser)
	if err != nil {
		// Handle error.
	}

	fmt.Printf("Values: %v\n", values)
	// Output: Values: map[Age:30 Username:srathi]
}
