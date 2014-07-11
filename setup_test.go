package multiconfig

import (
	"testing"
	"os"
	"flag"
	"reflect"
)

func bulkReplaceEnv(env map[string]string) map[string]string {
	old := make(map[string]string)
	if env == nil {
		return old
	}

	for k, v := range env {
		old[k] = os.Getenv(k)
		err := os.Setenv(k, v)
		if err != nil {
			panic(err)
		}
	}

	return old
}

func TestSetupInto(t *testing.T) {
	type basics struct {
		One string `default:"default"`
		Two bool `default:"true"`
		ThirdField int `default:"4"`
	}

	type envChange struct {
		One string `default:"def" env:"OTHER"`
	}

	tests := []struct {
		Obj interface{}
		Args []string
		Environment map[string]string
		Base string
		Output interface{}
	}{
		{&basics{}, nil, nil, "base", &basics{"default", true, 4}},
		{&basics{}, []string{"-one=not"}, nil, "base", &basics{"not", true, 4}},
		{&basics{}, []string{"--third-field", "8"}, nil, "base", &basics{"default", true, 8}},
		{&basics{}, nil, map[string]string{"BASE_THIRD_FIELD": "11"}, "base", &basics{"default", true, 11}},
		{&basics{}, []string{"--one", "two"}, map[string]string{"BASE_ONE": "three"}, "base", &basics{"two", true, 4}},
		{&envChange{}, nil, nil, "base", &envChange{"def"}},
		{&envChange{}, nil, map[string]string{"BASE_OTHER": "four"}, "base", &envChange{"four"}},
	}

	for _, test := range tests {
		oldEnv := bulkReplaceEnv(test.Environment)

		set := flag.NewFlagSet("", flag.ContinueOnError)
		err := SetupInto(test.Obj, test.Base, set)
		if err != nil {
			t.Errorf("SetupInto(%#v, %#v, set) returned unexpected error %v", test.Obj, test.Base, err)
		} else {
			err := set.Parse(test.Args)
			if err != nil {
				t.Fatalf("Unexpected error calling set.Parse for output %#v: %v", test.Output, err)
			}

			if !reflect.DeepEqual(test.Obj, test.Output) {
				t.Errorf("SetupInto(obj, %#v, set) = %#v, but wanted %#v", test.Base, test.Obj, test.Output)
			}
		}

		bulkReplaceEnv(oldEnv)
	}
}
