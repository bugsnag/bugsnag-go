package bugsnag

import (
	"github.com/bugsnag/bugsnag-go/errors"
	"reflect"
	"testing"
	"unsafe"
)

type _account struct {
	Id   string
	Name string
	Plan struct {
		Premium bool
	}
	secret string
}

type _broken struct {
	Me   *_broken
	Data string
}

var account = _account{}
var notifier = NewNotifier(Configuration{})

func TestMetaDataAdd(t *testing.T) {
	m := MetaData{
		"one": {
			"key":      "value",
			"override": false,
		}}

	m.Add("one", "override", true)
	m.Add("one", "new", "key")
	m.Add("new", "tab", account)

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
	account.Id = "test"
	account.secret = "hush"

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
			"password": "[REDACTED]",
			"array": []interface{}{map[string]interface{}{
				"creditcard": "[REDACTED]",
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
				"Id":   "test",
				"Name": "test",
				"Plan": map[string]interface{}{
					"Premium": false,
				},
			},
		},
	}) {
		t.Errorf("metadata.Sanitize didn't work: %#v", n)
	}

}

func ExampleMetaData() {
	notifier.Notify(errors.Errorf("hi world"),
		MetaData{"Account": {
			"id":      account.Id,
			"name":    account.Name,
			"paying?": account.Plan.Premium,
		}})
}
