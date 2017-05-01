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

func TestEnvMap(t *testing.T) {

	prefix := "TEST_"
	delimiter := "$"

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
