package jsondiff

import (
	"encoding/json"
	"log"
	"reflect"
)

type CompareJSONOptions struct {
	AddedKey   string `json:"added_key"`
	RemovedKey string `json:"removed_key"`
	ChangedKey string `json:"changed_key"`
}

func DefaultCompareJSONOptions() CompareJSONOptions {
	return CompareJSONOptions{
		AddedKey:   "added",
		RemovedKey: "removed",
		ChangedKey: "changed",
	}
}

func CompareJSON(a, b []byte, opts *CompareJSONOptions) (string, error) {
	var m1 map[string]interface{}
	var m2 map[string]interface{}

	err := json.Unmarshal(a, &m1)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(b, &m2)
	if err != nil {
		return "", err
	}

	result := make(map[string]interface{})

	added := make(map[string]interface{})
	removed := make(map[string]interface{})
	changed := make(map[string]interface{})

	addedKey := "added"
	if opts != nil && opts.AddedKey != "" {
		addedKey = opts.AddedKey
	}

	removedKey := "removed"
	if opts != nil && opts.RemovedKey != "" {
		removedKey = opts.RemovedKey
	}

	changedKey := "changed"
	if opts != nil && opts.ChangedKey != "" {
		changedKey = opts.ChangedKey
	}

	for k, v1 := range m1 {
		if v2, ok := m2[k]; ok {
			if !reflect.DeepEqual(v1, v2) {
				switch v1.(type) {
				case map[string]interface{}:
					bytev1, err := json.Marshal(v1)
					if err != nil {
						log.Printf("%v", err)
					}
					bytev2, err := json.Marshal(v2)
					if err != nil {
						log.Printf("%v", err)
					}
					subResult, err := CompareJSON(bytev1, bytev2, opts)
					if err == nil {
						if len(subResult) > 0 {
							var subMap map[string]interface{}
							err = json.Unmarshal([]byte(subResult), &subMap)
							if err == nil {
								result[k] = subMap
							} else {
								log.Printf("%v", err)
							}
						}
					}
				default:
					oldVal, _ := json.Marshal(v1)
					newVal, _ := json.Marshal(v2)
					changed[k] = map[string]interface{}{
						"old": json.RawMessage(oldVal),
						"new": json.RawMessage(newVal),
					}
				}
			}
		} else {
			removed[k] = v1
		}
	}

	for k, v2 := range m2 {
		if _, ok := m1[k]; !ok {
			added[k] = v2
		}
	}

	if len(removed) > 0 {
		result[removedKey] = removed
	}

	if len(added) > 0 {
		result[addedKey] = added
	}

	if len(changed) > 0 {
		result[changedKey] = changed
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}
