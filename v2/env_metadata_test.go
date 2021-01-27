package bugsnag

import "testing"

func TestParseMetadataKeypath(t *testing.T) {
	type output struct {
		keypath string
		err     string
	}
	cases := map[string]output{
		"":                                {"", "No metadata prefix found"},
		"BUGSNAG_METADATA_":               {"", "No metadata prefix found"},
		"BUGSNAG_METADATA_key":            {"key", ""},
		"BUGSNAG_METADATA_device_foo":     {"device_foo", ""},
		"BUGSNAG_METADATA_device_foo_two": {"device_foo_two", ""},
	}

	for input, expected := range cases {
		keypath, err := parseMetadataKeypath(input)
		if len(expected.err) > 0 && (err == nil || err.Error() != expected.err) {
			t.Errorf("expected error with message '%s', got '%v'", expected.err, err)
		}
		if expected.keypath != keypath {
			t.Errorf("expected keypath '%s', got '%s'", expected.keypath, keypath)
		}
	}
}

func TestLoadEnvMetadata(t *testing.T) {
	cases := map[string]envMetadata{
		"":                                     {"", "", ""},
		"BUGSNAG_METADATA_Orange=tomato_paste": {"custom", "Orange", "tomato_paste"},
		"BUGSNAG_METADATA_true_orange=tomato_paste":             {"true", "orange", "tomato_paste"},
		"BUGSNAG_METADATA_color_Orange=tomato_paste":            {"color", "Orange", "tomato_paste"},
		"BUGSNAG_METADATA_color_Orange_hue=tomato_paste":        {"color", "Orange_hue", "tomato_paste"},
		"BUGSNAG_METADATA_crayonColor_Magenta=tomato_paste":     {"crayonColor", "Magenta", "tomato_paste"},
		"BUGSNAG_METADATA_crayonColor_Magenta_hue=tomato_paste": {"crayonColor", "Magenta_hue", "tomato_paste"},
	}

	for input, expected := range cases {
		metadata := loadEnvMetadata([]string{input})

		if len(expected.tab) == 0 {
			for _, m := range metadata {
				t.Errorf("erroneously added a value for '%s' to tab '%s':'%s'", input, m.tab, m.key)
			}
		} else {
			if len(metadata) != 1 {
				t.Fatalf("wrong number of metadata elements: %d %v", len(metadata), metadata)
			}
			m := metadata[0]
			if m.tab != expected.tab {
				t.Errorf("wrong tab '%s'", expected.tab)
				continue
			}
			if m.key != expected.key {
				t.Errorf("wrong key '%s'", expected.key)
				continue
			}
			if m.value != expected.value {
				t.Errorf("incorrect value added to keypath: '%s'", m.value)
			}
		}
	}
}
