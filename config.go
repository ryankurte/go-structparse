package configurer

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func ParseConfig(delimiter, envPrefix string, c interface{}) {
	// Parse overrides
	InfillConfig(delimiter, envPrefix, c)
}

// LoadConfig loads application configuration from the provided yml file
// Then infills configuration variables from the environment
func LoadConfig(filename, delimiter, envPrefix string, c interface{}) error {

	// Load configuration file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Unmarshal from yaml
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}

	// Parse overrides
	ParseConfig(delimiter, envPrefix, c)

	return nil
}
