package utils

import "encoding/json"

func MergeStructsToMap(a, b interface{}) map[string]interface{} {
	mapA := structToMap(a)
	mapB := structToMap(b)

	// Merge maps
	for key, value := range mapB {
		mapA[key] = value
	}

	return mapA
}

func structToMap(data interface{}) map[string]interface{} {
	jsonData, _ := json.Marshal(data)
	var result map[string]interface{}
	json.Unmarshal(jsonData, &result)
	return result
}
