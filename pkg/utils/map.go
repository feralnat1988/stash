package utils

import (
	"strings"
)

// NestedMap is a map that supports nested keys.
// It is expected that the nested maps are of type map[string]interface{}
type NestedMap map[string]interface{}

func (m NestedMap) Get(key string) (interface{}, bool) {
	fields := strings.Split(key, ".")

	current := m

	for _, f := range fields[:len(fields)-1] {
		v, found := current[f]
		if !found {
			return nil, false
		}

		current, _ = v.(map[string]interface{})
		if current == nil {
			return nil, false
		}
	}

	ret, found := current[fields[len(fields)-1]]
	return ret, found
}

func (m NestedMap) Set(key string, value interface{}) {
	fields := strings.Split(key, ".")

	current := m

	for _, f := range fields[:len(fields)-1] {
		v, ok := current[f].(map[string]interface{})
		if !ok {
			v = make(map[string]interface{})
			current[f] = v
		}

		current = v
	}

	current[fields[len(fields)-1]] = value
}
