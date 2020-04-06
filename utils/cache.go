package utils

import "fmt"

func InterfacesToStrings(values ...interface{}) []string {
	valuesStr := make([]string, 0)
	for _, v := range values {
		valuesStr = append(valuesStr, fmt.Sprint(v))
	}
	return valuesStr
}
