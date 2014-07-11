package multiconfig

import (
	"testing"
)

func TestEnvName(t *testing.T) {
	tests := []struct {
		Base, Name, Output string
	}{
		{"base", "FieldName15", "BASE_FIELDNAME15"},
		{"prefixHere", "underscores_in_name", "PREFIXHERE_UNDERSCORES_IN_NAME"},
		{"extra_", "_UnderscoresHere", "EXTRA_UNDERSCORESHERE"},
		{"base", "dash-es", "BASE_DASH_ES"},
	}

	for _, test := range tests {
		got := EnvName(test.Base, test.Name)
		if got != test.Output {
			t.Errorf("EnvName(%#v, %#v) = %#v, but wanted %#v", test.Base, test.Name, got, test.Output)
		}
	}
}

func TestFlagName(t *testing.T) {
	tests := []struct {
		Name, Output string
	}{
		{"FieldName15", "field-name-15"},
		{"Field", "field"},
		{"ComboFieldName", "combo-field-name"},
		{"With_Underscores", "with-underscores"},
	}

	for _, test := range tests {
		got := FlagName(test.Name)
		if got != test.Output {
			t.Errorf("FlagName(%#v) = %#v, but wanted %#v", test.Name, got, test.Output)
		}
	}
}
