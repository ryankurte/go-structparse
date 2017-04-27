package configurer

import (
	"fmt"
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

func NoTestConfigInfill(t *testing.T) {

	t.Run("Infill strings in config structure", func(t *testing.T) {
		prefix := "TEST_"
		delimiter := "$"

		c := FakeConfig{
			KeyOne:   fmt.Sprintf("%s%s", delimiter, keyOneName),
			Map:      make(map[string]string),
			Embedded: Embedded{fmt.Sprintf("%s%s", delimiter, keyThreeName)},
		}
		c.Map[keyTwoIndex] = fmt.Sprintf("%s%s", delimiter, keyTwoName)

		os.Setenv(fmt.Sprintf("%s%s", prefix, keyOneName), keyOneValue)
		os.Setenv(fmt.Sprintf("%s%s", prefix, keyTwoName), keyTwoValue)
		os.Setenv(fmt.Sprintf("%s%s", prefix, keyThreeName), keyThreeValue)

		InfillConfig(delimiter, prefix, &c)

		assert.EqualValues(t, keyOneValue, c.KeyOne)
		assert.EqualValues(t, keyTwoValue, c.Map[keyTwoIndex])
		assert.EqualValues(t, keyThreeValue, c.Embedded.KeyThree)
	})

}
