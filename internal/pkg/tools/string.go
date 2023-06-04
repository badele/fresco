package tools

import (
	"strings"
)

func ConvertToQuotedSlice(value string) []string {
	arr := strings.Split(value, ",")
	for i := range arr {
		arr[i] = strings.TrimSpace(arr[i])
	}

	return arr
}

func UniqStringSlice(sslices ...[]string) []string {
	uniqueMap := map[string]bool{}
	for _, sslice := range sslices {
		for _, value := range sslice {
			uniqueMap[value] = true
		}
	}

	result := make([]string, 0, len(uniqueMap))
	for key := range uniqueMap {
		result = append(result, key)
	}

	return result
}
