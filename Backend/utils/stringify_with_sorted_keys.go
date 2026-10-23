package utils

import "encoding/json"

func StringifyWithSortedKeys(obj interface{}) (string, error) {
	// json.Marshal guarantees that map keys will be sorted when outputting JSON.
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
