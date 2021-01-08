package main

import (
	"testing"

	"github.com/bugsnag/bugsnag-go/v2"
)

func TestFormatMetadataKey(t *testing.T) {
	cases := map[string]string{
		"":                "",
		"final_level":     "final level",
		"Ecco_theDolphin": "Ecco theDolphin",
		"dot.delimiter":   "dot.delimiter",
	}

	for input, expected := range cases {
		output := formatMetadataKey(input)
		if output != expected {
			t.Errorf("expected '%s' got '%s'", expected, output)
		}
	}
}

func TestParseMetadataKeypath(t *testing.T) {
	type output struct {
		keypath string
		err     string
	}
	cases := map[string]output{
		"":                                {"", "No metadata prefix found"},
		"BUGSNAG_METADATA_":               {"", "No metadata prefix found"},
		"BUGSNAG_METADATA_key":            {"key", ""},
		"BUGSNAG_METADATA_device.foo":     {"device.foo", ""},
		"BUGSNAG_METADATA_device.foo.two": {"device.foo.two", ""},
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

func TestAddMetadata(t *testing.T) {
	type output struct {
		tab string
		key string
	}
	cases := map[string]output{
		"":                         {"", ""},
		"Orange":                   {"custom", "Orange"},
		"true_orange":              {"custom", "true orange"},
		"color.Orange":             {"color", "Orange"},
		"color.Orange_hue":         {"color", "Orange hue"},
		"crayon.color.Magenta":     {"crayon", "color.Magenta"},
		"crayon.color.Magenta_hue": {"crayon", "color.Magenta hue"},
	}

	for input, expected := range cases {
		event := &bugsnag.Event{MetaData: make(bugsnag.MetaData)}
		addMetadata(event, input, "tomato")
		if len(expected.tab) == 0 {
			for tab, pairs := range event.MetaData {
				for key, value := range pairs {
					if value == "tomato" {
						t.Errorf("erroneously added a value for '%s' to tab '%s':'%s'", input, tab, key)
					}
				}
			}
		} else {
			pairs, ok := event.MetaData[expected.tab]
			if !ok {
				t.Errorf("no value added to tab '%s'", expected.tab)
				continue
			}
			value, ok := pairs[expected.key]
			if !ok {
				t.Errorf("no value present for key '%s'", expected.key)
				continue
			}
			if value != "tomato" {
				t.Errorf("incorrect value added to keypath: '%s'", value)
			}
		}
	}
}
