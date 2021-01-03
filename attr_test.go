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

func TestGetValue(t *testing.T) {
	want := user.Username
	got, err := GetValue(user, "Username")
	require.Nil(t, err)
	require.Equal(t, want, got, "Username mismatch")

	wantErr := ErrUnexportedField
	_, gotErr := GetValue(user, "password")
	require.Equal(t, wantErr, gotErr, "Able to get an unexported field value")
}

func ExampleGetValue() {
	testUser := User{Username: "srathi", password: "secret", Age: 30}

	value, err := GetValue(testUser, "Age")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("Value of Age: %v", value)
	// Output: Value of Age: 30
}

func TestHas(t *testing.T) {
	for _, test := range []struct {
		attrName string
		want     bool
		errMsg   string
	}{
		{"Age", true, "Age not found"},
		{"NonExistant", false, "NonExistant field found"},
	} {
		got, err := Has(&user, test.attrName)
		require.Nil(t, err)
		require.Equal(t, test.want, got, test.errMsg)
	}
}

func ExampleHas() {
	testUser := User{Username: "srathi", password: "secret", Age: 30}

	ok, err := Has(&testUser, "Age")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("Age found: %v\n", ok)

	ok, err = Has(&testUser, "ABC")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("ABC found: %v\n", ok)
	// Output:
	// Age found: true
	// ABC found: false
}

func TestSetValue(t *testing.T) {
	for _, test := range []struct {
		attrName string
		newValue interface{}
		byPtr    bool
		wantErr  error
		errMsg   string
	}{
		{"Username", "new-srathi", false, ErrNotPtr, "Able to set fields on a struct by value"},
		{"password", "new-password", true, ErrUnexportedField, "Able to set a private field of a struct"},
		{"Age", 40.5, true, ErrMismatchValue, "Able to set float value to an int field"},
		{"Age", 40, true, nil, "New value not set in a struct"},
		{"Age", 30, true, nil, "New value not set in a struct"},
	} {
		var gotErr error
		if test.byPtr {
			gotErr = SetValue(&user, test.attrName, test.newValue)
		} else {
			gotErr = SetValue(user, test.attrName, test.newValue)
		}

		require.Equal(t, test.wantErr, gotErr, test.errMsg)
	}
}

func ExampleSetValue() {
	testUser := User{Username: "srathi", password: "secret", Age: 30}

	err := SetValue(&testUser, "password", "new-secret")
	fmt.Printf("Error while setting a private field: %v\n", err)

	err = SetValue(testUser, "Username", "new-username")
	fmt.Printf("Error while passing struct by value: %v\n", err)

	err = SetValue(&testUser, "Username", 100)
	fmt.Printf("Error while setting 100 in username: %v\n", err)

	err = SetValue(&testUser, "Username", "new-username")
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

func TestNames(t *testing.T) {
	// Only public fields are returned.
	want := []string{"Username", "Age"}
	got, err := Names(&user)
	require.Nil(t, err)
	require.Equal(t, want, got, "Struct field list is not correct")
}

func ExampleNames() {
	testUser := User{Username: "srathi", password: "secret", Age: 30}

	fields, err := Names(&testUser)
	if err != nil {
		// Handle error.
	}
	fmt.Printf("Field names: %v", fields)
	// Output: Field names: [Username Age]
}

func TestValues(t *testing.T) {
	// Only the value of the public fields are returned.
	want := map[string]interface{}{"Username": "srathi", "Age": 30}
	got, err := Values(&user)
	require.Nil(t, err)
	require.Equal(t, want, got, "Struct field values are not correct")
}

func ExampleValues() {
	testUser := User{Username: "srathi", password: "secret", Age: 30}

	values, err := Values(&testUser)
	if err != nil {
		// Handle error.
	}

	fmt.Printf("Values: %v\n", values)
	// Output: Values: map[Age:30 Username:srathi]
}

func TestGetTag(t *testing.T) {
	for _, test := range []struct {
		attrName  string
		tagName   string
		wantValue string
		errMsg    string
	}{
		{"Age", "meta", "important", "meta tag value for 'Age' is not correct"},
		{"Username", "json", "username", "json tag value for 'Username' is not correct"},
		{"Age", "db", "", "json tag value for 'Age' is not correct"},
	} {
		gotValue, err := GetTag(&user, test.attrName, test.tagName)
		require.Nil(t, err)
		require.Equal(t, test.wantValue, gotValue, test.errMsg)
	}

	wantErr := ErrUnexportedField
	_, gotErr := GetTag(&user, "password", "json")
	require.Equal(t, wantErr, gotErr, "Able to get tag value of a private field")
}

func ExampleGetTag() {
	// type User struct {
	// 	Username string `json:"username" db:"uname"`
	// 	password string `json:"password" db:"pw"`
	// 	Age      int    `json:"age" meta:"important"`
	// }
	testUser := User{Username: "srathi", password: "secret", Age: 30}
	tag, err := GetTag(&testUser, "Username", "db")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("db tag value: %s\n", tag)
	// Output: db tag value: uname
}

func TestTags(t *testing.T) {
	for _, test := range []struct {
		tagName string
		want    map[string]string
		errMsg  string
	}{
		{"json", map[string]string{"Username": "username", "Age": "age"},
			"json tag values are not correct"},
		{"db", map[string]string{"Username": "uname", "Age": ""},
			"db tag values are not correct"},
	} {
		got, err := Tags(&user, test.tagName)
		require.Nil(t, err)
		require.Equal(t, test.want, got, test.errMsg)
	}
}

func ExampleTags() {
	// type User struct {
	// 	Username string `json:"username" db:"uname"`
	// 	password string `json:"password" db:"pw"`
	// 	Age      int    `json:"age"  meta:"important"`
	// }
	testUser := User{Username: "srathi", password: "secret", Age: 30}

	values, err := Tags(&testUser, "json")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("json tag values: %v\n", values)

	values, err = Tags(&testUser, "db")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("db tag values: %v\n", values)
	// Output:
	// json tag values: map[Age:age Username:username]
	// db tag values: map[Age: Username:uname]
}

func TestGetKind(t *testing.T) {
	for _, test := range []struct {
		attrName string
		kindStr  string
	}{
		{"Username", "string"},
		{"Age", "int"},
	} {
		want := test.kindStr
		got, err := GetKind(&user, test.attrName)
		require.Nil(t, err)
		require.Equal(t, want, got, "Field 'Kind' mismatch")
	}

	wantErr := ErrNoField
	_, gotErr := GetKind(user, "ABC")
	require.Equal(t, wantErr, gotErr, "Able to get a non-existent field 'Kind'")
}

func ExampleGetKind() {
	testUser := User{Username: "srathi", password: "secret", Age: 30}

	kind, err := GetKind(testUser, "Age")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("Kind of Age: %s\n", kind)

	kind, err = GetKind(testUser, "Username")
	if err != nil {
		// Handle error.
	}
	fmt.Printf("Kind of Username: %s\n", kind)
	// Output:
	// Kind of Age: int
	// Kind of Username: string
}

func TestKinds(t *testing.T) {
	// Only public fields are returned.
	want := map[string]string{"Username": "string", "Age": "int"}
	got, err := Kinds(&user)
	require.Nil(t, err)
	require.Equal(t, want, got, "Struct field 'kind' map is not correct")
}

func ExampleKinds() {
	testUser := User{Username: "srathi", password: "secret", Age: 30}

	kinds, err := Kinds(&testUser)
	if err != nil {
		// Handle error.
	}
	fmt.Printf("Field kinds: %v", kinds)
	// Output: Field kinds: map[Age:int Username:string]
}
