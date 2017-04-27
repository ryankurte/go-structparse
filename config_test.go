package configurer

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type HostExample struct {
	Name string
	Port string
}

type ConfigExample struct {
	Name string
	Host HostExample
}

func TestConfigLoading(t *testing.T) {

	t.Run("Load and infill config file", func(t *testing.T) {
		delimiter := "$"
		prefix := "TEST_"

		os.Setenv(fmt.Sprintf("%s%s", prefix, "NAME"), "APP_NAME")
		os.Setenv(fmt.Sprintf("%s%s", prefix, "HOSTNAME"), "localhost")
		os.Setenv(fmt.Sprintf("%s%s", prefix, "PORT"), "9009")

		c := ConfigExample{}
		err := LoadConfig("./example.yml", delimiter, prefix, &c)
		assert.Nil(t, err)

		assert.EqualValues(t, "APP_NAME", c.Name)
		assert.EqualValues(t, "localhost", c.Host.Name)
		assert.EqualValues(t, "9009", c.Host.Port)

		log.Printf("Config: %+v", c)

	})

}
