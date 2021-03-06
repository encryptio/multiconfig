package multiconfig

import (
	"flag"
	"os"
	"reflect"
	"testing"
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
		One        string
		Two        bool `default:"true"`
		ThirdField int  `default:"4"`
	}

	type envChange struct {
		One string `default:"def" env:"OTHER"`
	}

	tests := []struct {
		Obj         interface{}
		Args        []string
		Environment map[string]string
		Base        string
		Output      interface{}
	}{
		{&basics{}, nil, nil, "base", &basics{"", true, 4}},
		{&basics{}, []string{"-one=not"}, nil, "base", &basics{"not", true, 4}},
		{&basics{}, []string{"--third-field", "8"}, nil, "base", &basics{"", true, 8}},
		{&basics{}, nil, map[string]string{"BASE_THIRD_FIELD": "11"}, "base", &basics{"", true, 11}},
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

func TestSetupIntoHelp(t *testing.T) {
	type helper struct {
		One string `help:"help for one"`
		Two string // no help
	}

	obj := &helper{}

	set := flag.NewFlagSet("", flag.ContinueOnError)
	err := SetupInto(obj, "base", set)
	if err != nil {
		t.Fatalf("SetupInto(helper{}, \"base\", set) returned unexpected error %v", err)
	}
	set.Parse(nil)

	checked := 0
	set.VisitAll(func(f *flag.Flag) {
		if f.Name == "one" {
			if f.Usage != "help for one" {
				t.Errorf("Help for One field is %#v, wanted %#v", f.Usage, "help for one")
			}
			checked++
		} else if f.Name == "two" {
			if f.Usage != "" {
				t.Errorf("Help for Two field is %#v, wanted %#v", f.Usage, "")
			}
			checked++
		}
	})
	if checked != 2 {
		t.Errorf("Didn't visit all flags")
	}
}
