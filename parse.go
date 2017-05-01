package structparse

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// NewEnvironmentMapper creates an environment mapping parser
// This parses a string looking for a delimiter indicating that the value should be loaded from the environment
func NewEnvironmentMapper(delimiter, prefix string) Parse {
	return func(line string) string {
		if !strings.HasPrefix(line, delimiter) {
			return line
		}
		key := fmt.Sprintf("%s%s", prefix, strings.Replace(line, delimiter, "", -1))
		value := os.Getenv(key)
		//log.Printf("Parsing: '%s' Key: '%s' Value: '%s'", line, key, value)
		return value
	}
}

// Parse is a string parsing function to be called when strings are found
type Parse func(in string) string

// Strings reflects over a structure and calls Parse when strings are located
func Strings(parse Parse, obj interface{}) {
	parseStringsRecursive(parse, reflect.ValueOf(obj))
}

// Internal recursive parsing function
func parseStringsRecursive(parse Parse, val reflect.Value) reflect.Value {
	switch val.Kind() {
	case reflect.Ptr:
		o := val.Elem()
		if o.IsValid() {
			parseStringsRecursive(parse, o)
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			res := parseStringsRecursive(parse, val.Field(i))
			if res != reflect.ValueOf(nil) {
				val.Field(i).Set(res)
			}
		}
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			res := parseStringsRecursive(parse, val.Index(i))
			if res != reflect.ValueOf(nil) {
				val.Index(i).Set(res)
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			mapVal := val.MapIndex(k)
			res := parseStringsRecursive(parse, mapVal)
			if res != reflect.ValueOf(nil) {
				val.SetMapIndex(k, res)
			}
		}
	case reflect.String:
		value := parse(val.String())
		return reflect.ValueOf(value)
	}

	return reflect.ValueOf(nil)
}
