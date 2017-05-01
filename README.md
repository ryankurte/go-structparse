# go-structparse

[![Documentation](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/ryankurte/ons)

A (deep) recursive field parser for golang. Basically just a reflection wrapper that finds elements in a struct or map and runs the provided function over them when found.
This is useful as a sort of "find and replace" for fields in a structure.

It was built to allow overriding of configuration variables in a config file with environmental variables, so that you can have one clear source of truth for configuration that explicitly defines any external requirements.

It is intended that this will be extended with further field parsing / better recursion as required.

## Usage

An `EnvironmentMapper` is provided that will parse strings looking for a delimiter, and when the delimiter is found replace the value with the environmental variable defined by the prefix and field value.
For example, with a delimiter of `$` and a prefix of `APP_` the value `$NAME` will be replaced with the envronmental variale `APP_NAME`.

To implement custom parsers, check out the godoc.


## Example

### Config file

``` yaml
---

# strings beginning with $ will be parsed as environmental variables with the provided prefix
name: $NAME
host: 
  name: $HOSTNAME
  port: $PORT

```

### Code
``` go
import (
    "io/ioutil"
    "gopkg.in/yaml.v2"

    "github.com/ryankurte/go-structparse"
)

type Config struct {
    Name string
    Host struct {
        Name string
        Port string
    }
}

func Main() {
    c := Config{}

    // Load (yaml) configuration file
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }

    // Unmarshal from yaml to config struct
    err = yaml.Unmarshal(data, c)
    if err != nil {
        return err
    }

    // Parse struct fields
    delimiter := "$"
    prefix := "APP_"
    structparse.Strings(structparse.NewEnvironmentMapper(delimiter, prefix), &c)

    // Do something with the parsed structure
    log.Printf("Config: %+v", c)

}

```

---

If you have any questions, comments, or suggestions, feel free to open an issue or a pull request.
