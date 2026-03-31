package utils

import (
	"net/url"
	"strconv"
	"strings"
)

// Int
func GetInt(q url.Values, key string, def int) int {
	val := q.Get(key)
	if v, err := strconv.Atoi(val); err == nil && v > 0 {
		return v
	}
	return def
}

func GetString(q url.Values, key, def string) string {
	val := strings.TrimSpace(q.Get(key))
	if val == "" {
		return def
	}
	return val
}

func GetBool(q url.Values, key string, def bool) bool {
	val := q.Get(key)
	if val == "" {
		return def
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return def
	}
	return b
}

// pointer values
func GetBoolPtr(q url.Values, key string) *bool {
	val := q.Get(key)
	if val == "" {
		return nil
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return nil
	}
	return &b
}

func GetIntPtr(q url.Values, key string) *int {
	val := q.Get(key)
	if val == "" {
		return nil
	}
	v, err := strconv.Atoi(val)
	if err != nil || v <= 0 {
		return nil
	}
	return &v
}

func GetStringPtr(q url.Values, key string) *string {
	val := strings.TrimSpace(q.Get(key))
	if val == "" {
		return nil
	}
	return &val
}
