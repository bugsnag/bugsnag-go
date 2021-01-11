package bugsnag

import (
	"fmt"
	"testing"
)

func TestParsePairs(t *testing.T) {
	type output struct {
		key, value string
		err error
	}

	cases := map[string]output{
		"":{"", "", fmt.Errorf("Not a '='-delimited key pair")},
		"key=value":{"key", "value", nil},
		"key=value=bar":{"key", "value=bar", nil},
		"something":{"", "", fmt.Errorf("Not a '='-delimited key pair")},
	}
	for input, expected := range cases {
		key, value, err := parseEnvironmentPair(input)
		if expected.err != nil && (err == nil || err.Error() != expected.err.Error()) {
			t.Errorf("expected error '%v', got '%v'", expected.err, err)
		}
		if key != expected.key || value != expected.value {
			t.Errorf("expected pair '%s'='%s', got '%s'='%s'", expected.key, expected.value, key, value)
		}
	}
}
