package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// Format {index} to args
func Format(s string, args ...interface{}) string {
	result := s
	for i, arg := range args {
		placeholder := "{" + strconv.Itoa(i) + "}"
		argString := fmt.Sprintf("%v", arg)
		result = strings.Replace(result, placeholder, argString, -1)
	}
	return result
}

// GetDefault allows the data param to be nil
func GetDefault(data any, def string) string {
	if data == nil {
		return def
	}
	return data.(string)
}
