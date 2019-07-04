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
		{name: "Gin", got: rv.Gin, exp: ""},
		{name: "Martini", got: rv.Martini, exp: ""},
		{name: "Negroni", got: rv.Negroni, exp: ""},
		{name: "Revel", got: rv.Revel, exp: ""},
	} {
		if tc.got != tc.exp {
			t.Errorf("expected pristine '%s' runtime version to be '%s' but was '%s'", tc.name, tc.exp, tc.got)
		}
	}
}

func TestModifiedRuntimeVersions(t *testing.T) {
	versions = nil // reset global variable
	rv := GetRuntimeVersions()
	AddVersion("Gin", "1.2.1")
	AddVersion("Martini", "1.0.0")
	AddVersion("Negroni", "1.0.2")
	AddVersion("Revel", "0.20.1")
	for _, tc := range []struct{ name, got, exp string }{
		{name: "Go", got: rv.Go, exp: runtime.Version()},
		{name: "Gin", got: rv.Gin, exp: "1.2.1"},
		{name: "Martini", got: rv.Martini, exp: "1.0.0"},
		{name: "Negroni", got: rv.Negroni, exp: "1.0.2"},
		{name: "Revel", got: rv.Revel, exp: "0.20.1"},
	} {
		if tc.got != tc.exp {
			t.Errorf("expected modified '%s' runtime version to be '%s' but was '%s'", tc.name, tc.exp, tc.got)
		}
	}

}
