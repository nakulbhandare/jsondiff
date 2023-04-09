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
	"changed": {
		"name": {
			"new": "Jane",
			"old": "John"
		}
	}
}
```

Here are some inputes and outputs

### Input 1

```
Json 1:
{
	"username": "John",
	"address": {
		"city": "New York",
		"country": "USA"
	},
	"gender": "male"
}

Json 2:

{
	"username": "John",
    "lastname": "wille"
	"age": 30,
	"address": {
		"city": "California",
		"country": "USA",
        "zipcode": 90210

	},
	"gender": "male"
}
```

This will output the following JSON string:

```
{
	"added": {
		"age": 30,
		"lastname": "wille"
	},
	"address": {
		"added": {
			"zipcode": 90210
		},
		"changed": {
			"city": {
				"new": "California",
				"old": "New York"
			}
		}
	}
}
```
### Input 2

```
Json 1:
{
	"name": "Alice",
	"age": 35,
	"address": {
	  "city": "Anytown",
	  "state": "CA",
	  "zip": "12345"
	},
	"preferences": {
	  "color": {
		"primary": "blue",
		"secondary": "green"
	  },
	  "food": "pizza"
	}
} 

Json 2:
{
	"name": "Alice",
	"age": 35,
	"address": {
		"city": "Anytown",
		"latitude": 37.7749,
		"longitude": -71.0589,
		"country": "USA"
	},
	"preferences": {
		"color": {
			"primary": "yellow",
			"secondary": "orange",
			"tertiary": "blue",
			"quaternary": "red"
		},
		"food": "pizza",
		"drink": "soda",
		"dessert": "cake"
	},
	"gender": "male",
	"profession": "engineer"
}
```

This will output the following JSON string:

```
{
	"added": {
		"gender": "male",
		"profession": "engineer"
	},
	"address": {
		"added": {
			"country": "USA",
			"latitude": 37.7749,
			"longitude": -71.0589
		},
		"removed": {
			"state": "CA",
			"zip": "12345"
		}
	},
	"preferences": {
		"added": {
			"dessert": "cake",
			"drink": "soda"
		},
		"color": {
			"added": {
				"quaternary": "red",
				"tertiary": "blue"
			},
			"changed": {
				"primary": {
					"new": "yellow",
					"old": "blue"
				},
				"secondary": {
					"new": "orange",
					"old": "green"
				}
			}
		}
	}
}
```
The library compares two JSON objects and collects the differences at each level. The differences are presented in a single object with properties for added, removed, and changed fields. Recursion is used to handle nested objects. This provides a concise and organized representation of the differences between the two JSON objects.

# Contributing

If you find a bug or have an idea for a new feature, feel free to open an issue or submit a pull request. All contributions are welcome!
