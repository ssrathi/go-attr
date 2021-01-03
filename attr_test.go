package attr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type User struct {
	Username string `json:"username" db:"uname"`
	Age      int    `json:"age" meta:"important"`
	password string
}

var user User = User{"srathi", 30, "my_secret_123"}

func TestGetField(t *testing.T) {
	want := user.Username
	got, err := GetField(user, "Username")
	require.Nil(t, err)
	require.Equal(t, want, got, "Username mismatch")

	wantErr := ErrUnexportedField
	_, gotErr := GetField(user, "password")
	require.Equal(t, wantErr, gotErr, "Able to get an unexported field value")
}

func ExampleGetField() {
	testUser := User{Username: "srathi", password: "secret", Age: 30}

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
	require.Equal(t, want, got, "Age not found")

	want = false
	got, err = HasField(&user, "NonExistant")
	require.Nil(t, err)
	require.Equal(t, want, got, "NonExistant field found")
}

func ExampleHasField() {
	testUser := User{Username: "srathi", password: "secret", Age: 30}

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
	require.Equal(t, wantErr, gotErr, "Able to set fields on a struct by value")

	wantErr = ErrUnexportedField
	gotErr = SetField(&user, "password", "new-password")
	require.Equal(t, wantErr, gotErr, "Able to set a private field of a struct")

	wantErr = ErrMismatchValue
	gotErr = SetField(&user, "Age", 40.5)
	require.Equal(t, wantErr, gotErr, "Able to set float value to an int field")

	oldVal := user.Age
	err := SetField(&user, "Age", 40)
	require.Nil(t, err)
	require.Equal(t, 40, user.Age, "New value not set in a struct")

	err = SetField(&user, "Age", oldVal)
	require.Nil(t, err)
	require.Equal(t, user.Age, oldVal, "New value not set in a struct")
}

func ExampleSetField() {
	testUser := User{Username: "srathi", password: "secret", Age: 30}

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
	require.Equal(t, want, got, "Struct field list is not correct")
}

func ExampleFieldNames() {
	testUser := User{Username: "srathi", password: "secret", Age: 30}

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
	require.Equal(t, want, got, "Struct field values are not correct")
}

func ExampleFieldValues() {
	testUser := User{Username: "srathi", password: "secret", Age: 30}

	values, err := FieldValues(&testUser)
	if err != nil {
		// Handle error.
	}

	fmt.Printf("Values: %v\n", values)
	// Output: Values: map[Age:30 Username:srathi]
}

func TestGetFieldTag(t *testing.T) {
	want := "important"
	got, err := GetFieldTag(&user, "Age", "meta")
	require.Nil(t, err)
	require.Equal(t, want, got, "meta tag value for 'Age' is not correct")

	want = "username"
	got, err = GetFieldTag(&user, "Username", "json")
	require.Nil(t, err)
	require.Equal(t, want, got, "json tag value for 'Username' is not correct")

	want = ""
	got, err = GetFieldTag(&user, "Age", "db")
	require.Nil(t, err)
	require.Equal(t, want, got, "json tag value for 'Age' is not correct")

	wantErr := ErrUnexportedField
	_, gotErr := GetFieldTag(&user, "password", "json")
	require.Equal(t, wantErr, gotErr, "Able to get tag value of a private field")
}

func ExampleGetFieldTag() {
	// type User struct {
	// 	Username string `json:"username" db:"uname"`
	// 	password string `json:"password" db:"pw"`
	// 	Age      int    `json:"age" meta:"important"`
	// }
	testUser := User{Username: "srathi", password: "secret", Age: 30}
	tag, err := GetFieldTag(&testUser, "Username", "db")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("db tag value: %s\n", tag)
	// Output: db tag value: uname
}

func TestTagValues(t *testing.T) {
	want := map[string]string{"Username": "username", "Age": "age"}
	got, err := TagValues(&user, "json")
	require.Nil(t, err)
	require.Equal(t, want, got, "json tag values are not correct")

	want = map[string]string{"Username": "uname", "Age": ""}
	got, err = TagValues(&user, "db")
	require.Nil(t, err)
	require.Equal(t, want, got, "db tag values are not correct")
}

func ExampleTagValues() {
	// type User struct {
	// 	Username string `json:"username" db:"uname"`
	// 	password string `json:"password" db:"pw"`
	// 	Age      int    `json:"age"  meta:"important"`
	// }
	testUser := User{Username: "srathi", password: "secret", Age: 30}

	values, err := TagValues(&testUser, "json")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("json tag values: %v\n", values)

	values, err = TagValues(&testUser, "db")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("db tag values: %v\n", values)
	// Output:
	// json tag values: map[Age:age Username:username]
	// db tag values: map[Age: Username:uname]
}

func TestGetFieldKind(t *testing.T) {
	for _, a := range []struct {
		attrName string
		kindStr  string
	}{
		{"Username", "string"}, {"Age", "int"},
	} {
		want := a.kindStr
		got, err := GetFieldKind(&user, a.attrName)
		require.Nil(t, err)
		require.Equal(t, want, got, "Field 'Kind' mismatch")
	}

	wantErr := ErrNoField
	_, gotErr := GetField(user, "ABC")
	require.Equal(t, wantErr, gotErr, "Able to get a non-existant field 'Kind'")
}

func ExampleGetFieldKind() {
	testUser := User{Username: "srathi", password: "secret", Age: 30}

	kind, err := GetFieldKind(testUser, "Age")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("Kind of Age: %s\n", kind)

	kind, err = GetFieldKind(testUser, "Username")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("Kind of Username: %s\n", kind)
	// Output:
	// Kind of Age: int
	// Kind of Username: string
}

func TestFieldKinds(t *testing.T) {
	// Only public fields are returned.
	want := map[string]string{"Username": "string", "Age": "int"}
	got, err := FieldKinds(&user)
	require.Nil(t, err)
	require.Equal(t, want, got, "Struct field 'kind' map is not correct")
}

func ExampleFieldKinds() {
	testUser := User{Username: "srathi", password: "secret", Age: 30}

	kinds, err := FieldKinds(&testUser)
	if err != nil {
		// Handle error.
	}
	fmt.Printf("Field kinds: %v", kinds)
	// Output: Field kinds: map[Age:int Username:string]
}
