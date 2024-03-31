package repositories

import (
	"strings"

	"golang.org/x/exp/maps"
)

func FieldsAndArgsFromMap(itemArgs map[string]any, delimiter string) (string, string) {
	fields := maps.Keys(itemArgs)
	namedArgs := ""
	for i, field := range fields {
		if i == len(fields)-1 {
			namedArgs += "@" + field
			break
		}
		namedArgs += "@" + field + ","
	}
	return strings.Join(fields, delimiter), namedArgs
}

func FieldsAndArgsFromSlice(itemArgs []string, delimiter string) (string, string) {
	namedArgs := ""
	for i, field := range itemArgs {
		if i == len(itemArgs)-1 {
			namedArgs += "@" + field
			break
		}
		namedArgs += "@" + field + ","
	}
	return strings.Join(itemArgs, delimiter), namedArgs
}

func Fields(fields []string, delimiter string) string {
	return strings.Join(fields, delimiter)
}
