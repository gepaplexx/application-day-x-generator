package util

import "fmt"

type Value struct {
	Val any
}

func (v Value) String() string {
	switch v.Val.(type) {
	case string:
		return fmt.Sprintf("%s", v.Val.(string))
	case bool:
		return fmt.Sprintf("%t", v.Val)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v.Val)
	case float32, float64:
		return fmt.Sprintf("%f", v.Val)
	}

	return NOT_VALID
}

func (v Value) Equal(other any) bool {
	switch v.Val.(type) {
	case string:
		return v.Val.(string) == other.(string)
	case bool:
		return v.Val.(bool) == other.(bool)
	}

	return false
}

func (v Value) IsEmpty() bool {
	switch v.Val.(type) {
	case string:
		return len(v.Val.(string)) > 0
	}

	return v.Val == nil
}
