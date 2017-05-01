package structparse

import (
	"fmt"
	"os"
	"strings"
)

type EnvironmentMapper struct {
	delimiter string
	prefix    string
}

// NewEnvironmentMapper creates an environment mapping parser
// This parses a string looking for a delimiter indicating that the value should be loaded from the environment
func NewEnvironmentMapper(delimiter, prefix string) *EnvironmentMapper {
	return &EnvironmentMapper{delimiter, prefix}
}

func (em *EnvironmentMapper) ParseString(line string) string {
	if !strings.HasPrefix(line, em.delimiter) {
		return line
	}
	key := fmt.Sprintf("%s%s", em.prefix, strings.Replace(line, em.delimiter, "", -1))
	value := os.Getenv(key)
	//log.Printf("Parsing: '%s' Key: '%s' Value: '%s'", line, key, value)
	return value
}
