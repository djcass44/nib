package env

import (
	"os"
	"strings"
)

type Source func(key string) string

func GetFirst(key string, src Source) string {
	values := Get(key, src)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func GetLast(key string, src Source) string {
	values := Get(key, src)
	if len(values) == 0 {
		return ""
	}
	return values[len(values)-1]
}

func Get(key string, src Source) []string {
	if src == nil {
		src = os.Getenv
	}
	return strings.Split(src(key), ":")
}
