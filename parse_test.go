/**
 * go-structparse
 * Copyright 2017 Ryan Kurte
 */

package structparse

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"
)

type FakeEnvMapper struct {
	Match   string
	Replace string
}

func (fem *FakeEnvMapper) ParseString(line string) string {
	if line == fem.Match {
		return fem.Replace
	}
	return line
}

func TestParsing(t *testing.T) {

	prefix := "TEST_"
	delimiter := "$"

	t.Run("Handles struct fields", func(t *testing.T) {
		fem := FakeEnvMapper{"TEST", "REPLACED"}
		c := struct{ Test string }{fem.Match}

		Strings(&fem, &c)

		assert.EqualValues(t, fem.Replace, c.Test)
	})

	t.Run("Handles map fields", func(t *testing.T) {
		fem := FakeEnvMapper{"TEST", "REPLACED"}
		c := make(map[string]string)
		c["test"] = fem.Match

		Strings(&fem, &c)

		assert.EqualValues(t, fem.Replace, c["test"])
	})

	t.Run("Handles embedded structs", func(t *testing.T) {
		fem := FakeEnvMapper{"TEST", "REPLACED"}
		c := struct{ Fake struct{ Test string } }{struct{ Test string }{fem.Match}}

		Strings(&fem, &c)

		assert.EqualValues(t, fem.Replace, c.Fake.Test)
	})

	t.Run("Handles embedded maps", func(t *testing.T) {
		fem := FakeEnvMapper{"TEST", "REPLACED"}
		c := struct{ Fake map[string]string }{make(map[string]string)}
		c.Fake["test"] = fem.Match

		Strings(&fem, &c)

		assert.EqualValues(t, fem.Replace, c.Fake["test"])
	})

	t.Run("Handles recursive maps", func(t *testing.T) {
		fem := FakeEnvMapper{"TEST", "REPLACED"}
		c := struct{ Fake map[string]map[string]string }{make(map[string]map[string]string)}
		c.Fake["test1"] = make(map[string]string)
		c.Fake["test1"]["test2"] = fem.Match

		Strings(&fem, &c)

		assert.EqualValues(t, fem.Replace, c.Fake["test1"]["test2"])
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
