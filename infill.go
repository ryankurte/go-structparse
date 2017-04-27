package configurer

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/oleiade/reflections"
)

// Parses a config struct and loads
func infillConfig(delimiter, prefix string, config interface{}) interface{} {

	original := reflect.ValueOf(config)

	copy := reflect.New(original.Type()).Elem()
	translateRecursive(delimiter, prefix, copy, original)

	// Remove the reflection wrapper
	return copy.Interface()
}

func v2(delimiter, prefix string, config interface{}) error {

	log.Printf("Object: %+v", config)

	fields, err := reflections.FieldsDeep(config)
	if err != nil {
		return err
	}

	log.Printf("Fields: %+v", fields)

	for _, fieldName := range fields {
		kind, err := reflections.GetFieldKind(config, fieldName)
		if err != nil {
			return err
		}

		field, err := reflections.GetField(config, fieldName)
		if err != nil {
			return err
		}

		switch kind {
		case reflect.String:
			value := mapString(delimiter, prefix, field.(string))
			reflections.SetField(config, fieldName, value)

		case reflect.Struct:
			v2(delimiter, prefix, reflect.ValueOf(field))

		case reflect.Map:
			v2(delimiter, prefix, field)
		}
	}

	return nil
}

func v2recurse(delimiter, prefix string, obj interface{}) error {
	val := reflect.ValueOf(obj)
	kind := val.Kind()
	t := val.Type()

	log.Printf("Object: %+v Kind: %+v", obj, kind)

	switch kind {
	case reflect.Ptr:
		if val.IsValid() {
			v2recurse(delimiter, prefix, val.Elem())
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			log.Printf("Field: %s %+v", t.Field(i).Name, reflect.ValueOf(field).Interface())
			//v2recurse(delimiter, prefix, field.Interface())
		}
	case reflect.String:
		value := mapString(delimiter, prefix, val.String())
		val.Elem().SetString(value)
	}

	return nil
}

// Parses a string looking for a delimiter indicating that the value should be loaded from the environment
func mapString(delimiter, prefix, line string) string {
	if !strings.HasPrefix(line, delimiter) {
		return line
	}

	key := fmt.Sprintf("%s%s", prefix, strings.Replace(line, delimiter, "", -1))
	value := os.Getenv(key)

	return value
}

// https://gist.github.com/hvoecking/10772475
func translateRecursive(delimiter, prefix string, copy, original reflect.Value) {
	switch original.Kind() {
	// The first cases handle nested structures and translate them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		translateRecursive(delimiter, prefix, copy.Elem(), originalValue)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		translateRecursive(delimiter, prefix, copyValue, originalValue)
		copy.Set(copyValue)

	// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			translateRecursive(delimiter, prefix, copy.Field(i), original.Field(i))
		}

	// If it is a slice we create a new slice and translate each element
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			translateRecursive(delimiter, prefix, copy.Index(i), original.Index(i))
		}

	// If it is a map we create a new map and translate each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			translateRecursive(delimiter, prefix, copyValue, originalValue)
			copy.SetMapIndex(key, copyValue)
		}

	// Otherwise we cannot traverse anywhere so this finishes the the recursion

	// If it is a string translate it (yay finally we're doing what we came for)
	case reflect.String:
		translatedString := mapString(delimiter, prefix, original.Interface().(string))
		copy.SetString(translatedString)

	// And everything else will simply be taken from the original
	default:
		copy.Set(original)
	}

}
