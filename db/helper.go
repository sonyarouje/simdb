package db

import (
	"os"
	"strings"
	"fmt"
	"strconv"
	"errors"
)

func createDirIfNotExist(dir string) error {

	if _,err:= os.Stat(dir); err==nil {
		return nil
	}
	err:=os.Mkdir(dir, 0755)
	return err
}


func mergeToExisting(array []interface{}, entity interface{}) ([]interface{}, error) {
	array=append(array, entity)
	return array, nil
}


// getNestedValue fetch nested value from node
func getNestedValue(input interface{}, node string) (interface{}, error) {
	pp := strings.Split(node, ".")
	for _, n := range pp {
		if isIndex(n) {
			// find slice/array
			if arr, ok := input.([]interface{}); ok {
				indx, err := getIndex(n)
				if err != nil {
					return input, err
				}
				arrLen := len(arr)
				if arrLen == 0 ||
					indx > arrLen-1 {
					return empty, errors.New("empty array")
				}
				input = arr[indx]
			}
		} else {
			// find in map
			validNode := false
			if mp, ok := input.(map[string]interface{}); ok {
				input, ok = mp[n]
				validNode = ok
			}

			// find in group data
			if mp, ok := input.(map[string][]interface{}); ok {
				input, ok = mp[n]
				validNode = ok
			}

			if !validNode {
				return empty, fmt.Errorf("invalid node name %s", n)
			}
		}
	}

	return input, nil
}

func getIndex(in string) (int, error) {
	if !isIndex(in) {
		return -1, fmt.Errorf("invalid index")
	}
	is := strings.TrimLeft(in, "[")
	is = strings.TrimRight(is, "]")
	oint, err := strconv.Atoi(is)
	if err != nil {
		return -1, err
	}
	return oint, nil
}

func isIndex(in string) bool {
	return strings.HasPrefix(in, "[") && strings.HasSuffix(in, "]")
}

// toFloat64 converts interface{} value to float64 if value is numeric else return false
func toFloat64(v interface{}) (float64, bool) {
	var f float64
	flag := true
	// as Go convert the json Numeric value to float64
	switch u := v.(type) {
	case int:
		f = float64(u)
	case int8:
		f = float64(u)
	case int16:
		f = float64(u)
	case int32:
		f = float64(u)
	case int64:
		f = float64(u)
	case float32:
		f = float64(u)
	case float64:
		f = u
	default:
		flag = false
	}
	return f, flag
}

// length return length of strings/array/map
func length(v interface{}) (int, error) {
	if val, ok := v.(string); ok {
		return len(val), nil
	} else if val, ok := v.([]interface{}); ok {
		return len(val), nil
	} else if val, ok := v.(map[string]interface{}); ok {
		return len(val), nil
	}
	return -1, errors.New("invalid type for length")
}
