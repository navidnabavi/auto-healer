package autoheal

import (
	"fmt"
	"os"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func checkErrorForIntegerEnv(err error, fieldName string) {
	if err == nil {
		return
	}
	fmt.Printf("field %s must be integer.\r\n", fieldName)
	panic(err)
}

func valueInList(value string, list map[string]string) bool {
	for k := range list {
		if value == k {
			return true
		}
	}
	return false
}
