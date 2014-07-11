// multiconfig reads configuration information from command line flags and environment variables.
//
// To use it, create a struct to contain your configuration variables and pass a pointer to it
// to ReadConfig.
package multiconfig

import (
	"reflect"
	"errors"
	"strconv"
	"time"
	"os"
	"flag"
)

var (
	ErrNotPointerToStruct = errors.New("conf object was not a pointer to a struct")
	ErrBadType = errors.New("conf object has a field with an unsupported type")
)

type ErrorFlag struct {
	FlagName string
	InitialVal string
	Err error
}

func (e ErrorFlag) Error() string {
	return "Can't parse "+e.FlagName+" value "+e.InitialVal+": "+e.Err.Error()
}

// Setup calls SetupInto(..., flag.CommandLine).
func Setup(conf interface{}, base string) error {
	return SetupInto(conf, base, flag.CommandLine)
}

// SetupInto reflects on the pointer to a struct passed and calls set.*Var to set up flags it finds.
// Note that you still need to call set.Parse to actually read configuration values.
//
// The default value given to flag.*Var is the value of the environment variable for the flag if
// set, otherwise the default given in the field tag, otherwise the zero value for the type.
//
// You may set tags on the struct fields to customize behavior:
//     type Conf struct {
//         Field1 string        `help:"Use with care"`
//         Field2 int           `env:"OVERRIDED_ENV_NAME"`
//         Field3 time.Duration `default:"3s"`
//         Field4 bool          `flag:"four"`
//     }
//
// If the "help" tag exists, it is passed as the help argument to flag.Var. If the "env" tag exists,
// it overrides the default environment variable name (see EnvName.) If the "flag" tag exists, it
// overrides the default flag name (see FlagName.)
//
// The supported types are the same as in the flag package: bool, time.Duration, float64, int,
// int64, string, uint, and uint64.
func SetupInto(conf interface{}, base string, set *flag.FlagSet) error {
	val := reflect.ValueOf(conf)

	if val.Kind() != reflect.Ptr {
		return ErrNotPointerToStruct
	}

	val = val.Elem()

	if val.Kind() != reflect.Struct {
		return ErrNotPointerToStruct
	}

	for i := 0; i < val.NumField(); i++ {
		f := val.Type().Field(i)

		defVal := f.Tag.Get("default")

		flagName := FlagName(f.Name)

		envName := f.Tag.Get("env")
		if envName == "" {
			envName = EnvName(base, flagName)
		}
		envVal := os.Getenv(envName)

		initialVal := envVal
		if envVal == "" {
			initialVal = defVal
		}

		help := f.Tag.Get("help")

		val := val.Field(i).Addr().Interface()
		switch val := val.(type) {
		case *bool:
			parsed, err := strconv.ParseBool(initialVal)
			if err != nil { return ErrorFlag{flagName, initialVal, err} }
			set.BoolVar(val, flagName, parsed, help)

		case *time.Duration:
			parsed, err := time.ParseDuration(initialVal)
			if err != nil { return ErrorFlag{flagName, initialVal, err} }
			set.DurationVar(val, flagName, parsed, help)

		case *int:
			parsed, err := strconv.ParseInt(initialVal, 0, 0)
			if err != nil { return ErrorFlag{flagName, initialVal, err} }
			set.IntVar(val, flagName, int(parsed), help)

		case *int64:
			parsed, err := strconv.ParseInt(initialVal, 0, 64)
			if err != nil { return ErrorFlag{flagName, initialVal, err} }
			set.Int64Var(val, flagName, parsed, help)

		case *uint:
			parsed, err := strconv.ParseUint(initialVal, 0, 0)
			if err != nil { return ErrorFlag{flagName, initialVal, err} }
			set.UintVar(val, flagName, uint(parsed), help)

		case *uint64:
			parsed, err := strconv.ParseUint(initialVal, 0, 64)
			if err != nil { return ErrorFlag{flagName, initialVal, err} }
			set.Uint64Var(val, flagName, uint64(parsed), help)

		case *string:
			set.StringVar(val, flagName, initialVal, help)

		case *float64:
			parsed, err := strconv.ParseFloat(initialVal, 64)
			if err != nil { return ErrorFlag{flagName, initialVal, err} }
			set.Float64Var(val, flagName, parsed, help)

		default:
			return ErrBadType
		}
	}

	return nil
}
