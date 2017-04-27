package configurer

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Embedded struct {
	KeyThree string
}

type FakeConfig struct {
	KeyOne   string
	Map      map[string]string
	Embedded Embedded
}

const (
	keyOneName    = "KEY_ONE"
	keyOneValue   = "KEY_ONE_VALUE"
	keyTwoIndex   = "KEY_TWO_INDEX"
	keyTwoName    = "KEY_TWO"
	keyTwoValue   = "KEY_TWO_VALUE"
	keyThreeName  = "KEY_THREE"
	keyThreeValue = "KEY_THREE_VALUE"
)

func TestEnvironmentInfill(t *testing.T) {

	t.Run("Infill strings in structure from the environment", func(t *testing.T) {
		prefix := "TEST_"
		delimiter := "$"

		c := FakeConfig{
			KeyOne:   fmt.Sprintf("%s%s", delimiter, keyOneName),
			Map:      make(map[string]string),
			Embedded: Embedded{fmt.Sprintf("%s%s", delimiter, keyThreeName)},
		}
		c.Map[keyTwoIndex] = fmt.Sprintf("%s%s", delimiter, keyTwoName)

		log.Printf("Config: %+v", c)

		os.Setenv(fmt.Sprintf("%s%s", prefix, keyOneName), keyOneValue)
		os.Setenv(fmt.Sprintf("%s%s", prefix, keyTwoName), keyTwoValue)
		os.Setenv(fmt.Sprintf("%s%s", prefix, keyThreeName), keyThreeValue)

		c2 := infillEnvironment(delimiter, prefix, &c)

		config2 := c2.(*FakeConfig)

		assert.EqualValues(t, keyOneValue, config2.KeyOne)
		assert.EqualValues(t, keyTwoValue, config2.Map[keyTwoIndex])
		assert.EqualValues(t, keyThreeValue, config2.Embedded.KeyThree)

		log.Printf("Config2: %+v", config2)

	})

}
