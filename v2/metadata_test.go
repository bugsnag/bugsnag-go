package bugsnag

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/bugsnag/bugsnag-go/v2/errors"
)

type _account struct {
	ID   string
	Name string
	Plan struct {
		Premium bool
	}
	Password      string
	secret        string
	Email         string `json:"email"`
	EmptyEmail    string `json:"emptyemail,omitempty"`
	NotEmptyEmail string `json:"not_empty_email,omitempty"`
}

type _broken struct {
	Me   *_broken
	Data string
}

var account = _account{}
var notifier = New(Configuration{})

func TestMetaDataAdd(t *testing.T) {
	m := MetaData{
		"one": {
			"key":      "value",
			"override": false,
		}}

	m.Add("one", "override", true)
	m.Add("one", "new", "key")
	m.Add("new", "tab", account)

	m.AddStruct("lol", "not really a struct")
	m.AddStruct("account", account)

	if !reflect.DeepEqual(m, MetaData{
		"one": {
			"key":      "value",
			"override": true,
			"new":      "key",
		},
		"new": {
			"tab": account,
		},
		"Extra data": {
			"lol": "not really a struct",
		},
		"account": {
			"ID":   "",
			"Name": "",
			"Plan": map[string]interface{}{
				"Premium": false,
			},
			"Password": "",
			"email":    "",
		},
	}) {
		t.Errorf("metadata.Add didn't work: %#v", m)
	}
}

func TestMetaDataUpdate(t *testing.T) {

	m := MetaData{
		"one": {
			"key":      "value",
			"override": false,
		}}

	m.Update(MetaData{
		"one": {
			"override": true,
			"new":      "key",
		},
		"new": {
			"tab": account,
		},
	})

	if !reflect.DeepEqual(m, MetaData{
		"one": {
			"key":      "value",
			"override": true,
			"new":      "key",
		},
		"new": {
			"tab": account,
		},
	}) {
		t.Errorf("metadata.Update didn't work: %#v", m)
	}
}

func TestMetaDataSanitize(t *testing.T) {

	var broken = _broken{}
	broken.Me = &broken
	broken.Data = "ohai"
	account.Name = "test"
	account.ID = "test"
	account.secret = "hush"
	account.Email = "example@example.com"
	account.EmptyEmail = ""
	account.NotEmptyEmail = "not_empty_email@example.com"

	m := MetaData{
		"one": {
			"bool":     true,
			"int":      7,
			"float":    7.1,
			"complex":  complex(1, 1),
			"func":     func() {},
			"unsafe":   unsafe.Pointer(broken.Me),
			"string":   "string",
			"password": "secret",
			"array": []hash{{
				"creditcard": "1234567812345678",
				"broken":     broken,
			}},
			"broken":  broken,
			"account": account,
		},
	}

	n := m.sanitize([]string{"password", "creditcard"})

	if !reflect.DeepEqual(n, map[string]interface{}{
		"one": map[string]interface{}{
			"bool":     true,
			"int":      7,
			"float":    7.1,
			"complex":  "[complex128]",
			"string":   "string",
			"unsafe":   "[unsafe.Pointer]",
			"func":     "[func()]",
			"password": "[FILTERED]",
			"array": []interface{}{map[string]interface{}{
				"creditcard": "[FILTERED]",
				"broken": map[string]interface{}{
					"Me":   "[RECURSION]",
					"Data": "ohai",
				},
			}},
			"broken": map[string]interface{}{
				"Me":   "[RECURSION]",
				"Data": "ohai",
			},
			"account": map[string]interface{}{
				"ID":   "test",
				"Name": "test",
				"Plan": map[string]interface{}{
					"Premium": false,
				},
				"Password":        "[FILTERED]",
				"email":           "example@example.com",
				"not_empty_email": "not_empty_email@example.com",
			},
		},
	}) {
		t.Errorf("metadata.Sanitize didn't work: %#v", n)
	}

}

func TestSanitizerSanitize(t *testing.T) {
	var (
		nilPointer   *int
		nilInterface = interface{}(nil)
	)

	for n, tc := range []struct {
		input interface{}
		want  interface{}
	}{
		{nilPointer, "<nil>"},
		{nilInterface, "<nil>"},
	} {
		s := &Sanitizer{}
		gotValue := s.Sanitize(tc.input)

		if got, want := gotValue, tc.want; got != want {
			t.Errorf("[%d] got %v, want %v", n, got, want)
		}
	}
}

func ExampleMetaData() {
	notifier.Notify(errors.Errorf("hi world"),
		MetaData{"Account": {
			"id":      account.ID,
			"name":    account.Name,
			"paying?": account.Plan.Premium,
		}})
}
