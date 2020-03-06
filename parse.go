package structparse

import (
	"reflect"
)

// StringParser interface is called on any string elements
type StringParser interface {
	ParseString(in string) interface{}
}

// IntParser called on integer elements
type IntParser interface {
	ParseInt(in int64) interface{}
}

// FloatParser called on floating point elements
type FloatParser interface {
	ParseFloat(in float64) interface{}
}

// Parsers is a container for all possible parser interfaces
type Parsers struct {
	StringParser StringParser
	IntParser    IntParser
	FloatParser  FloatParser
}

// Strings reflects over a structure and calls Parse when strings are located
func Strings(parser StringParser, obj interface{}) {
	parsers := Parsers{
		StringParser: parser,
	}
	parseRecursive(parsers, reflect.ValueOf(obj))
}

// Parse an object using the provided parser
// The parser must implement one of more ObjectParser interfaces
func Parse(parser, obj interface{}) {
	parsers := Parsers{}

	if p, ok := parser.(StringParser); ok {
		parsers.StringParser = p
	}
	if p, ok := parser.(IntParser); ok {
		parsers.IntParser = p
	}
	if p, ok := parser.(FloatParser); ok {
		parsers.FloatParser = p
	}

	parseRecursive(parsers, reflect.ValueOf(obj))
}

// Internal recursive parsing function
func parseRecursive(parsers Parsers, val reflect.Value) reflect.Value {
	switch val.Kind() {
	case reflect.Ptr:
		o := val.Elem()
		if o.IsValid() {
			res := parseRecursive(parsers, o)
			if res != reflect.ValueOf(nil) {
				return res
			}
		}
	case reflect.Interface:
		res := parseRecursive(parsers, val.Elem())
		if res != reflect.ValueOf(nil) {
			return res
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			res := parseRecursive(parsers, val.Field(i))
			if res != reflect.ValueOf(nil) {
				if val.Field(i).CanSet() {
					val.Field(i).Set(res)
				}
			}
		}
		return val
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			res := parseRecursive(parsers, val.Index(i))
			if res != reflect.ValueOf(nil) {
				if val.Index(i).CanSet() {
					val.Index(i).Set(res)
				}
			}
		}
		return val
	case reflect.Map:
		for _, k := range val.MapKeys() {
			mapVal := val.MapIndex(k)
			res := parseRecursive(parsers, mapVal)
			if res != reflect.ValueOf(nil) {
				val.SetMapIndex(k, res)
			}
		}
		return val
	case reflect.String:
		if parsers.StringParser != nil && val.Type().AssignableTo(reflect.TypeOf("")) {
			value := parsers.StringParser.ParseString(val.String())
			return reflect.ValueOf(value)
		}
		return val
	case reflect.Int:
		if parsers.IntParser != nil {
			value := parsers.IntParser.ParseInt(val.Int())
			return reflect.ValueOf(value)
		}
		return val
	case reflect.Float64:
		if parsers.FloatParser != nil {
			value := parsers.FloatParser.ParseFloat(val.Float())
			return reflect.ValueOf(value)
		}
		return val
	default:
		return val
	}

	return reflect.ValueOf(nil)
}
