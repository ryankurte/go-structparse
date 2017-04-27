# go-parse

```
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

// Parse struct fields


```