package utils

import "os"

func GetEnv(key, defaults string) (value string) {
	value = defaults
	if v, ok := os.LookupEnv(key); ok {
		value = v
	}

	return
}
