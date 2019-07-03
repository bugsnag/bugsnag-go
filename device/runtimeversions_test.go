package device

import (
	"runtime"
	"testing"
)

func TestPristineRuntimeVersions(t *testing.T) {
	versions = nil // reset global variable
	rv := GetRuntimeVersions()
	for _, tc := range []struct{ name, got, exp string }{
		{name: "Go", got: rv.Go, exp: runtime.Version()},
		{name: "gin", got: rv.Gin, exp: ""},
		{name: "martini", got: rv.Martini, exp: ""},
		{name: "negroni", got: rv.Negroni, exp: ""},
		{name: "revel", got: rv.Revel, exp: ""},
	} {
		if tc.got != tc.exp {
			t.Errorf("expected pristine '%s' runtime version to be '%s' but was '%s'", tc.name, tc.exp, tc.got)
		}
	}
}

func TestModifiedRuntimeVersions(t *testing.T) {
	versions = nil // reset global variable
	rv := GetRuntimeVersions()
	AddVersion("gin", "1.2.1")
	AddVersion("martini", "1.0.0")
	AddVersion("negroni", "1.0.2")
	AddVersion("revel", "0.20.1")
	for _, tc := range []struct{ name, got, exp string }{
		{name: "Go", got: rv.Go, exp: runtime.Version()},
		{name: "gin", got: rv.Gin, exp: "1.2.1"},
		{name: "martini", got: rv.Martini, exp: "1.0.0"},
		{name: "negroni", got: rv.Negroni, exp: "1.0.2"},
		{name: "revel", got: rv.Revel, exp: "0.20.1"},
	} {
		if tc.got != tc.exp {
			t.Errorf("expected pristine '%s' runtime version to be '%s' but was '%s'", tc.name, tc.exp, tc.got)
		}
	}

}
