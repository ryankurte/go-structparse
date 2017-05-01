/**
 * go-structparse
 * Copyright 2017 Ryan Kurte
 */

package structparse

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"
)

const (
	prefix    = "TEST_"
	delimiter = "$"
)

func GetKey(i uint64) string {
	return fmt.Sprintf("KEY%d", i)
}

func GetValue(i uint64) string {
	return fmt.Sprintf("VALUE%d", i)
}

var keys = make(map[string]string)

func RegisterKey(i uint64) {
	name := fmt.Sprintf("%s%s", prefix, GetKey(i))
	keys[name] = GetValue(i)
}

func ClearKey(i uint64) {
	name := fmt.Sprintf("%s%s", prefix, GetKey(i))
	keys[name] = ""
}

func FakeEnvMapper(line string) string {
	key := fmt.Sprintf("%s%s", prefix, strings.Replace(line, delimiter, "", -1))
	val, ok := keys[key]
	if !ok {
		return "ERROR"
	}
	return val
}

func TestParsing(t *testing.T) {

	RegisterKey(1)

	t.Run("Handles struct fields", func(t *testing.T) {
		c := struct{ Test string }{GetKey(1)}

		Strings(FakeEnvMapper, &c)

		assert.EqualValues(t, GetValue(1), c.Test)
	})

	t.Run("Handles map fields", func(t *testing.T) {
		c := make(map[string]string)
		c["test"] = GetKey(1)

		Strings(FakeEnvMapper, &c)

		assert.EqualValues(t, GetValue(1), c["test"])
	})

	t.Run("Handles embedded structs", func(t *testing.T) {
		c := struct{ Fake struct{ Test string } }{struct{ Test string }{GetKey(1)}}

		Strings(FakeEnvMapper, &c)

		assert.EqualValues(t, GetValue(1), c.Fake.Test)
	})

	t.Run("Handles embedded maps", func(t *testing.T) {
		c := struct{ Fake map[string]string }{make(map[string]string)}
		c.Fake["test"] = GetKey(1)

		Strings(FakeEnvMapper, &c)

		assert.EqualValues(t, GetValue(1), c.Fake["test"])
	})

	t.Run("Handles recursive maps", func(t *testing.T) {
		c := struct{ Fake map[string]map[string]string }{make(map[string]map[string]string)}
		c.Fake["test1"] = make(map[string]string)
		c.Fake["test1"]["test2"] = GetKey(1)

		Strings(FakeEnvMapper, &c)

		assert.EqualValues(t, GetValue(1), c.Fake["test1"]["test2"])
	})

	t.Run("Infills strings in structure from environment", func(t *testing.T) {

		os.Setenv(fmt.Sprintf("%s%s", prefix, "NAME"), "APP_NAME")
		os.Setenv(fmt.Sprintf("%s%s", prefix, "HOSTNAME"), "localhost")
		os.Setenv(fmt.Sprintf("%s%s", prefix, "PORT"), "9009")

		c := struct {
			Name string
			Host struct {
				Name string
				Port string
			}
		}{}

		// Load configuration file
		data, err := ioutil.ReadFile("./example.yml")
		assert.Nil(t, err)

		// Unmarshal from yaml
		err = yaml.Unmarshal(data, &c)
		assert.Nil(t, err)

		Strings(NewEnvironmentMapper(delimiter, prefix), &c)

		assert.EqualValues(t, "APP_NAME", c.Name)
		assert.EqualValues(t, "localhost", c.Host.Name)
		assert.EqualValues(t, "9009", c.Host.Port)
	})

}
