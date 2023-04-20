package jsondiff

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
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
	var m1 interface{}
	var m2 interface{}

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

	switch m1.(type) {
	case []interface{}:
		switch m2.(type) {
		case []interface{}:
			m1Array := m1.([]interface{})
			m2Array := m2.([]interface{})
			if len(m1Array) != len(m2Array) {
				return "", errors.New("arrays have different lengths")
			}
			for i := 0; i < len(m1Array); i++ {
				m1Obj, ok1 := m1Array[i].(map[string]interface{})
				m2Obj, ok2 := m2Array[i].(map[string]interface{})
				if !ok1 || !ok2 {
					return "", errors.New("array elements are not objects")
				}
				objChanged := make(map[string]interface{})
				for k, v1 := range m1Obj {
					v2, ok := m2Obj[k]
					if !ok {
						removed[k] = v1
					} else if !reflect.DeepEqual(v1, v2) {
						objChanged[k] = []interface{}{v1, v2}
					}
				}
				for k, v2 := range m2Obj {
					_, ok := m1Obj[k]
					if !ok {
						added[k] = v2
					}
				}
				if len(objChanged) > 0 {
					changed[strconv.Itoa(i)] = objChanged
				}
			}
		default:
			removed = m1.(map[string]interface{})
			added = m2.(map[string]interface{})
		}
	case map[string]interface{}:
		switch m2.(type) {
		case map[string]interface{}:
			compareMaps(m1.(map[string]interface{}), m2.(map[string]interface{}), &added, &removed, &changed, opts)
		default:
			removed = m1.(map[string]interface{})
			added = m2.(map[string]interface{})
		}
	default:
		if reflect.DeepEqual(m1, m2) {
			return "{}", nil
		}
		removed = m1.(map[string]interface{})
		added = m2.(map[string]interface{})
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

func compareMaps(m1, m2 map[string]interface{}, added, removed, changed *map[string]interface{}, opts *CompareJSONOptions) {
	for k, v1 := range m1 {
		if v2, ok := m2[k]; ok {
			switch v1.(type) {
			case []interface{}:
				switch v2.(type) {
				case []interface{}:
					subAdded := make([]interface{}, 0)
					subRemoved := make([]interface{}, 0)
					subChanged := make([]map[string]interface{}, 0)
					compareArrays(v1.([]interface{}), v2.([]interface{}), &subAdded, &subRemoved, &subChanged, opts)
					if len(subAdded) > 0 {
						(*added)[k] = subAdded
					}
					if len(subRemoved) > 0 {
						(*removed)[k] = subRemoved
					}
					if len(subChanged) > 0 {
						(*changed)[k] = subChanged
					}
				default:
					oldVal, _ := json.Marshal(v1)
					newVal, _ := json.Marshal(v2)
					(*changed)[k] = map[string]interface{}{
						"old": json.RawMessage(oldVal),
						"new": json.RawMessage(newVal),
					}
				}
			case map[string]interface{}:
				switch v2.(type) {
				case map[string]interface{}:
					subAdded := make(map[string]interface{})
					subRemoved := make(map[string]interface{})
					subChanged := make(map[string]interface{})
					compareMaps(v1.(map[string]interface{}), v2.(map[string]interface{}), &subAdded, &subRemoved, &subChanged, opts)
					if len(subAdded) > 0 {
						(*added)[k] = subAdded
					}
					if len(subRemoved) > 0 {
						(*removed)[k] = subRemoved
					}
					if len(subChanged) > 0 {
						(*changed)[k] = subChanged
					}
				default:
					oldVal, _ := json.Marshal(v1)
					newVal, _ := json.Marshal(v2)
					(*changed)[k] = map[string]interface{}{
						"old": json.RawMessage(oldVal),
						"new": json.RawMessage(newVal),
					}
				}
			default:
				if !reflect.DeepEqual(v1, v2) {
					oldVal, _ := json.Marshal(v1)
					newVal, _ := json.Marshal(v2)
					(*changed)[k] = map[string]interface{}{
						"old": json.RawMessage(oldVal),
						"new": json.RawMessage(newVal),
					}
				}
			}
		} else {
			(*removed)[k] = v1
		}
	}

	for k, v2 := range m2 {
		if _, ok := m1[k]; !ok {
			(*added)[k] = v2
		}
	}
}

func compareArrays(a, b []interface{}, added, removed *[]interface{}, changed *[]map[string]interface{}, opts *CompareJSONOptions) {
	for i := 0; i < len(a) || i < len(b); i++ {
		if i >= len(b) {
			*removed = append(*removed, a[i])
			continue
		}
		if i >= len(a) {
			*added = append(*added, b[i])
			continue
		}
		switch aValue := a[i].(type) {
		case map[string]interface{}:
			if bMap, ok := b[i].(map[string]interface{}); ok {
				itemDiff := make(map[string]interface{})
				hasChanged := false
				addedFields := make(map[string]interface{})
				removedFields := make(map[string]interface{})
				for k, v := range aValue {
					if bv, ok := bMap[k]; !ok {
						removedFields[k] = v
						hasChanged = true
					} else if !reflect.DeepEqual(v, bv) {
						itemDiff[k] = map[string]interface{}{"old": v, "new": bv}
						hasChanged = true
					} else {
						itemDiff[k] = v
					}
				}
				for k, v := range bMap {
					if _, ok := aValue[k]; !ok {
						addedFields[k] = v
						hasChanged = true
					}
				}
				if len(addedFields) > 0 {
					addedKey := "added"
					if opts.AddedKey != "" {
						addedKey = opts.AddedKey
					}
					itemDiff[addedKey] = addedFields
				}
				if len(removedFields) > 0 {
					removedKey := "removed"
					if opts.RemovedKey != "" {
						removedKey = opts.RemovedKey
					}
					itemDiff[removedKey] = removedFields
				}
				if hasChanged {
					*changed = append(*changed, itemDiff)
				}
			} else {
				*removed = append(*removed, aValue)
				*added = append(*added, b[i])
			}
		case []interface{}:
			if bArray, ok := b[i].([]interface{}); ok {
				compareArrays(aValue, bArray, added, removed, changed, opts)
			} else {
				*removed = append(*removed, aValue)
				*added = append(*added, b[i])
			}
		default:
			if !reflect.DeepEqual(aValue, b[i]) {
				changedKey := "changed"
				if opts.ChangedKey != "" {
					changedKey = opts.ChangedKey
				}
				switch reflect.ValueOf(aValue).Kind() {
				case reflect.Map, reflect.Slice:
					*changed = append(*changed, map[string]interface{}{changedKey: []interface{}{aValue, b[i]}})
				default:
					*changed = append(*changed, map[string]interface{}{changedKey: map[string]interface{}{"old": aValue, "new": b[i]}})
				}
			}
		}
	}
}
