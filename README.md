# go-structparse

[![Documentation](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/ryankurte/go-structparse)
[![GitHub tag](https://img.shields.io/github/tag/ryankurte/go-structparse.svg)](https://github.com/ryankurte/go-structparse)
[![Build Status](https://travis-ci.org/ryankurte/go-structparse.svg?branch=master)](https://travis-ci.org/ryankurte/go-structparse)

A (deep) recursive field parser for golang. Basically just a reflection wrapper that finds elements in a struct or map and runs the provided function over them when found.
This is useful as a sort of "find and replace" for fields or types in a structure.

It was built to allow overriding of configuration variables in a config file with environmental variables, so that you can have one clear source of truth for configuration that explicitly defines any external requirements.

It is intended that this will be extended with further field / type parsing as required. If you have a need, or an idea, open an issue!

## Usage

Sick of passing 4 million command line arguments into your app? Done with configuring unknown / unspecified environmental variables to get your app to work? Well, this might be the project for you!

What if you could have one, human readable, source of truth, that also supported the loading of secrets from the environment in an EXPLICIT manner?

Do we have the solution for you! Try our new `EnvironmentMapper` which will parse strings looking for a delimiter, and when the delimiter is found replace the value with the environmental variable defined by the prefix and field value.
For example, if you configure the `EnvironmentMapper` with a delimiter of `$` and a prefix of `APP_`, the value `$NAME` in your configuration file will be replaced with the environmental variable `APP_NAME`.
Now you have a static, explicit configuration file, that supports injection of secrets and other per-instance data. How neat is that!

If you fancy implementing custom parsers, check out the [godocs](https://godoc.org/github.com/ryankurte/go-structparse) for more information.


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
