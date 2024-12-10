package maps

import "reflect"

func DeepCopy[K string|int|uint, V interface{}](original map[K]V) map[K]V {
	copy := make(map[K]V)
    for key, value := range original {
		valueIsMap := reflect.TypeOf(value).Kind() == reflect.Map
		if valueIsMap {
			originalVal, ok := any(value).(map[K]V)
			if ok {
				copy[key] = any(DeepCopy(originalVal)).(V)
			} else {
                panic("Unsupported map type for deep copy")
            }
		} else {
			copy[key] = value
		}
	}

	return copy
}