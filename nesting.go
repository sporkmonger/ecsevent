package ecsevent

import (
	"fmt"
	"strings"
)

// Nest converts a map from dotted notation to a fully nested representation.
func Nest(entry map[string]interface{}) map[string]interface{} {
	var ok bool
	newEntry := make(map[string]interface{})
	for key, value := range entry {
		segments := strings.Split(key, ".")
		currObj := newEntry
		nextObj := newEntry
		for i, segment := range segments {
			if i < len(segments)-1 {
				// internal node
				nextObj, ok = (currObj[segment]).(map[string]interface{})
				if !ok {
					nextObj = make(map[string]interface{})
					currObj[segment] = nextObj
				}
				currObj = nextObj
			} else {
				// leaf node
				currObj[segment] = value
			}
		}
	}
	return newEntry
}

// Unnest converts a map from a nested representation into flat dotted notation.
func Unnest(entry map[string]interface{}) map[string]interface{} {
	newEntry := make(map[string]interface{})
	for key, value := range entry {
		if mapValue, ok := value.(map[string]interface{}); ok {
			// Nested value
			unnestedMap := Unnest(mapValue)
			for suffixKey, value := range unnestedMap {
				newEntry[fmt.Sprintf("%s.%s", key, suffixKey)] = value
			}
		} else {
			// Scalar value or list
			newEntry[key] = value
		}
	}
	return newEntry
}
