package sliceutils

import (
	"errors"
	"fmt"
)

func Containts(values interface{}, param interface{}) (bool, error) {
	switch t := values.(type) {
	case []int:
		for _, v := range t {
			if v == param {
				return true, nil
			}
		}
	case []int32:
		for _, v := range t {
			if v == param {
				return true, nil
			}
		}
	case []int64:
		for _, v := range t {
			if v == param {
				return true, nil
			}
		}
	case []float32:
		for _, v := range t {
			if v == param {
				return true, nil
			}
		}
	case []float64:
		for _, v := range t {
			if v == param {
				return true, nil
			}
		}
	case []string:
		for _, v := range t {
			if v == param {
				return true, nil
			}
		}
	default:
		return false, errors.New("无匹配类型")
	}
	return false, nil
}

func RemoveElement2(values interface{}, param interface{}) interface{} {
	switch t := values.(type) {
	case []int32:
		for k, v := range t {
			if v == param {
				fmt.Println(v, k)
				t = append(t[:k], t[k+1:]...)
				return t
			}
		}
	case []string:
		for k, v := range t {
			if v == param {
				fmt.Println(v, k)
				t = append(t[:k], t[k+1:]...)
				return t
			}
		}
	case []int64:
		for k, v := range t {
			if v == param {
				fmt.Println(v, k)
				t = append(t[:k], t[k+1:]...)
				return t
			}
		}
	default:
		return nil
	}
	return errors.New("slice is nil")
}

func RemoveInt32(values []int32, params int32) []int32 {
	NewValues := make([]int32, 0)
	for k, v := range values {
		if v == params {
			NewValues = append(values[:k], values[k+1:]...)
			return NewValues
		}
	}
	fmt.Println("slice:", NewValues)
	return NewValues

}
