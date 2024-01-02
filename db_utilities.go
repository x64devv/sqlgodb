package main 

import (
	"fmt"
	"strings"
)

func insertStmtStr (size int) string {
	var str []string ;
	for i:= 0; i < size; i ++{
		str = append(str, "?");
	}
	return strings.Join(str, ",")
}

func updateStmtStr(slice []string)  string{
	var str []string;
	for _, key := range slice{
		str = append(str, fmt.Sprintf("%s = ?", key))
	}
	return strings.Join(str, ",");
}

func arrayToInterfaceArr[T any](array []T) []interface{}{
	var arrayInterface []interface{} 
	for _, val := range array{
		arrayInterface = append(arrayInterface, val)
	}
	return arrayInterface
}

func MapKeyValue(m map[string]string, saveKeys bool) []string {
	var slice []string
	for key, value := range m {
		if saveKeys {
			slice = append(slice, key)
			continue
		}
		slice = append(slice, value)
	}
	return slice
}