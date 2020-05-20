/**
 * go-structparse
 * Copyright 2017 Ryan Kurte
 */

package structparse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type FakeEnvMapper struct {
	Match   string
	Replace string
}

func (fem *FakeEnvMapper) ParseString(line string) interface{} {
	if line == fem.Match {
		return fem.Replace
	}
	return line
}
func TestParsing(t *testing.T) {

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
		c.Fake["test2"] = make(map[string]string)
		c.Fake["test1"]["test1"] = fem.Match
		c.Fake["test1"]["test2"] = "boop"
		c.Fake["test2"]["test1"] = "boop"

		Strings(&fem, &c)

		assert.EqualValues(t, fem.Replace, c.Fake["test1"]["test1"])
		assert.EqualValues(t, "boop", c.Fake["test1"]["test2"])
		assert.EqualValues(t, "boop", c.Fake["test2"]["test1"])
	})

	t.Run("Handles non-supported types", func(t *testing.T) {
		fem := FakeEnvMapper{"TEST", "REPLACED"}
		c := struct{ Fake bool }{false}
		Strings(&fem, &c)
	})

	t.Run("struct with time", func(t *testing.T) {
		fem := FakeEnvMapper{"TEST", "REPLACED"}
		type Foo struct {
			Bar  string
			Time time.Time
		}

		foo := Foo{Bar: "TEST"}

		Strings(&fem, &foo)

		assert.EqualValues(t, "REPLACED", foo.Bar)
	})

	t.Run("map", func(t *testing.T) {
		fem := FakeEnvMapper{"TEST", "REPLACED"}
		type foo struct {
			Name string
		}
		c := map[string]foo{"1": {"TEST"}, "2": {"TEST"}}

		Strings(&fem, &c)

		assert.EqualValues(t, "REPLACED", c["1"].Name)
		assert.EqualValues(t, "REPLACED", c["2"].Name)
	})

	t.Run("struct", func(t *testing.T) {
		fem := FakeEnvMapper{"TEST", "REPLACED"}
		type bar struct {
			Name string
		}
		type foo struct {
			Bar bar
		}

		c := foo{Bar: bar{"TEST"}}

		Strings(&fem, &c)

		assert.EqualValues(t, "REPLACED", c.Bar.Name)
	})

	t.Run("struct with pointer field", func(t *testing.T) {
		fem := FakeEnvMapper{"TEST", "REPLACED"}
		type bar struct {
			Name       string
			unexported string
		}
		type foo struct {
			Bar *bar
		}

		c := foo{Bar: &bar{Name: "TEST", unexported: "unexported"}}

		Strings(&fem, &c)

		assert.EqualValues(t, "REPLACED", c.Bar.Name)
		assert.EqualValues(t, "unexported", c.Bar.unexported)
	})
}
