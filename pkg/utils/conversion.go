package utils

import "time"

func NilOrString(value *string) string {
	if value != nil {
		return *value
	}
	return ""
}

func NilOrFloat64(value *float64) *float64 {
	return value
}

func NilOrInt(value *int) *int {
	return value
}

func NilOrBool(value *bool) bool {
	if value != nil {
		return *value
	}
	return false
}

func StringOrNil(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func Float64OrNil(value float64) *float64 {
	if value == 0 {
		return nil
	}
	return &value
}

// needed

func Ptr[T any](v T) *T { return &v }

func TimeToString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func ParseTime(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
