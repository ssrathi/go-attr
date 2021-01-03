package attr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type User struct {
	Username  string
	password  string
	FirstName string
}

var user User = User{
	Username:  "srathi",
	password:  "my_secret_123",
	FirstName: "Shyamsunder",
}

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
		Username:  "srathi",
		FirstName: "Shyamsunder",
	}

	value, err := GetField(testUser, "FirstName")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("Value of FirstName: %s", value)
	// Output: Value of FirstName: Shyamsunder
}

func TestHasField(t *testing.T) {
	want := true
	got, err := HasField(&user, "FirstName")
	require.Nil(t, err)
	require.Equal(t, got, want, "FirstName not found")

	want = false
	got, err = HasField(&user, "NonExistant")
	require.Nil(t, err)
	require.Equal(t, got, want, "NonExistant field found")
}

func ExampleHasField() {
	testUser := User{
		Username:  "srathi",
		FirstName: "Shyamsunder",
	}

	ok, err := HasField(&testUser, "FirstName")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("FirstName found: %v\n", ok)

	ok, err = HasField(&testUser, "ABC")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("ABC found: %v\n", ok)

	// Output:
	// FirstName found: true
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
	gotErr = SetField(&user, "FirstName", 10)
	require.Equal(t, gotErr, wantErr, "Able to set int value to a string field")

	err := SetField(&user, "FirstName", "new-firstname")
	require.Nil(t, err)
	require.Equal(t, user.FirstName, "new-firstname", "New value not set in a struct")
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
