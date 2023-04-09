# jsondiff

The `jsondiff` is a Go library for computing the differences between two JSON documents. It provides an easy way for comparing JSON documents and returning the differences in a human-readable as well as machine-readable format. This library is useful for a wide range of applications, including testing, debugging, and versioning of JSON data.


# Installation

To install the `jsondiff` package, run the following command:

`go get github.com/nakulbhandare/jsondiff`

# Usage

The `CompareJSON` function in the `jsondiff` package takes two JSON byte arrays as input and returns a JSON string that contains the differences between the two JSON objects. The function signature is as follows:

```
func CompareJSON(a, b []byte, opts *CompareJSONOptions) (string, error)
```

`CompareJSONOptions` is a struct that contains the options for the comparison. By default, the options are as follows:

```
{
    "added_key": "added",
    "removed_key": "removed",
    "changed_key": "changed"
}
```

These options can be customized by passing a `CompareJSONOptions` struct to the CompareJSON function.

Here's an examples of how to use the `CompareJSON` function:

## Example 1

```
package main

import (
    "fmt"
    "github.com/nakulbhandare/jsondiff"
)

func main() {
    a := []byte(`{"name": "John", "age": 30}`)
    b := []byte(`{"name": "Jane", "age": 30, "city": "New York"}`)

    diff, err := jsondiff.CompareJSON(a, b, nil)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(diff)
}
```

This will output the following JSON string:

```
{
    "added": {
        "city": "New York"
    },
    "removed": {
        "name": "John"
    },
    "changed": {
        "name": {
            "old": "\"John\"",
            "new": "\"Jane\""
        }
    }
}
```

# Contributing

If you find a bug or have an idea for a new feature, feel free to open an issue or submit a pull request. All contributions are welcome!
