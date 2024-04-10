package bugsnag

import (
	"encoding/json"
	stderrors "errors"
	"reflect"
	"testing"
	"time"
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

type _textMarshaller struct{}

func (_textMarshaller) MarshalText() ([]byte, error) {
	return []byte("marshalled text"), nil
}

type _testStringer struct{}

func (s _testStringer) String() string {
	return "something"
}

type _testError struct{}

func (s _testError) Error() string {
	return "errorstr"
}

type _testStruct struct {
	Name *_testStringer
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

func TestMetadataAddPointer(t *testing.T) {
	var pointer *_testStringer
	md := MetaData{}
	md.AddStruct("emptypointer", pointer)
	fullPointer := &_testStringer{}
	md.AddStruct("fullpointer", fullPointer)

	if !reflect.DeepEqual(md, MetaData{
		"Extra data": {
			"emptypointer": "<nil>",
			"fullpointer":  "something",
		},
	}) {
		t.Errorf("metadata.AddStruct didn't work: %#v", md)
	}
}

func TestMetadataAddNil(t *testing.T) {
	md := MetaData{}
	md.AddStruct("map", map[string]interface{}{
		"data": _testStruct{Name: nil},
	})

	var nilMap map[string]interface{}
	md.AddStruct("nilmap", nilMap)

	var nilError _testError
	md.AddStruct("error", nilError)

	var nilErrorPtr *_testError
	md.AddStruct("errorNilPtr", nilErrorPtr)

	var timeVar time.Time
	md.AddStruct("timeUnset", timeVar)

	var duration time.Duration
	md.AddStruct("durationUnset", duration)

	var marshalNilPtr *_textMarshaller
	md.AddStruct("marshalNilPtr", marshalNilPtr)

	var marshalFullPtr = &_textMarshaller{}
	md.AddStruct("marshalFullPtr", marshalFullPtr)

	if !reflect.DeepEqual(md, MetaData{
		"map": {
			"data": map[string]interface{}{
				"Name": "<nil>",
			},
		},
		"nilmap": map[string]interface{}{},
		"Extra data": {
			"error":          "errorstr",
			"errorNilPtr":    "<nil>",
			"timeUnset":      "0001-01-01T00:00:00Z",
			"durationUnset":  "0s",
			"marshalFullPtr": "marshalled text",
			"marshalNilPtr":  "<nil>",
		},
	}) {
		t.Errorf("metadata.AddStruct didn't work: %#v", md)
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
			"error":    stderrors.New("some error"),
			"time":     time.Date(2023, 12, 5, 23, 59, 59, 123456789, time.UTC),
			"duration": 105567462 * time.Millisecond,
			"text":     _textMarshaller{},
			"json":     json.RawMessage(`{"json_property": "json_value"}`),
			"bytes":    []byte(`lots of bytes`),
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
			"error":    "some error",
			"time":     "2023-12-05T23:59:59.123456789Z",
			"duration": "29h19m27.462s",
			"text":     "marshalled text",
			"json":     `{"json_property": "json_value"}`,
			"bytes":    "lots of bytes",
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
		s := &sanitizer{}
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
