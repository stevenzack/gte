package util

import (
	"errors"
	"reflect"
	"strconv"
)

//ReplaceFieldIND replace field if not default value
func ReplaceFieldIND(dst, replacement interface{}) error {
	target := reflect.ValueOf(dst)
	if target.Kind() != reflect.Ptr {
		return errors.New("dst is not pointer value")
	}
	target = target.Elem()
	v := reflect.ValueOf(replacement)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		field := v.Field(i)
		// continue if field name doesn't match
		if _, ok := target.Type().FieldByName(fieldName); !ok {
			continue
		}
		switch field.Type().Kind() {
		case reflect.Bool, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64, reflect.String, reflect.Float32, reflect.Float64, reflect.Map, reflect.Slice:
			if !field.IsZero() {
				target.FieldByName(fieldName).Set(field)
			}
		}
	}
	return nil
}

func MapOf(kvs ...interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if len(kvs) == 0 {
		return m, nil
	}
	for i := 0; i < len(kvs); i++ {
		key, ok := kvs[i].(string)
		if !ok {
			return nil, errors.New("mapOf() failed: key[" + strconv.Itoa(i) + "] is not a string type")
		}
		if i+1 >= len(kvs) {
			return nil, errors.New("mapOf() failed: key[" + strconv.Itoa(i) + "] has no value in arguments")
		}
		value := kvs[i+1]
		m[key] = value
		i++
	}
	return m, nil
}

