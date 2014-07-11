package multiconfig

import (
	"strings"
	"unicode"
)

// EnvName is the default mapping from field names to environment variables. It takes the output
// of FlagName, uppercases every character, replaces dashes with underscores, and prepends the
// uppercase of base+"_".
//
// For example:
//    EnvName("program", "FieldName15") -> "PROGRAM_FIELD_NAME_15"
func EnvName(base, varName string) string {
	ret := FlagName(varName)

	ret = strings.ToUpper(ret)
	ret = strings.Replace(ret, "-", "_", -1)
	ret = strings.ToUpper(base) + "_" + ret

	n := ""
	for n != ret {
		n = ret
		ret = strings.Replace(ret, "__", "_", -1)
	}

	return ret
}

// FlagName is the default mapping from from field names to flag names. It adds dashes between
// pairs of characters when the case switches from lower to uppercase and when unicode.IsLetter
// changes.
//
// Finally, it maps all letters to lowercase, replaces underscores with dashes, and removes
// repeated dashes.
//
// For example:
//     FlagName("FieldName15") -> "field-name-15"
func FlagName(varName string) string {
	ret := ""

	last := ' '
	for i, c := range varName {
		if i == 0 {
			// skip checks for the first character
		} else if unicode.IsLetter(last) != unicode.IsLetter(c) {
			ret += "-"
		} else if unicode.IsLetter(c) && unicode.IsLower(last) && unicode.IsUpper(c) {
			ret += "-"
		}

		ret += string(c)
		last = c
	}

	ret = strings.ToLower(ret)
	ret = strings.Replace(ret, "_", "-", -1)

	n := ""
	for n != ret {
		n = ret
		ret = strings.Replace(ret, "--", "-", -1)
	}

	return ret
}
