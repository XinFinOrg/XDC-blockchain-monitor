package notification

import (
	"fmt"
	"reflect"

	"github.com/liam-lai/xinfin-monitor/types"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func GetMessageForSlack(bc *types.Blockchain) string {
	v := reflect.ValueOf(bc)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var message string
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		fieldName := fieldType.Name

		// Check if the field has the slack:"ignore" tag
		if tag, ok := fieldType.Tag.Lookup("slack"); ok && tag == "ignore" {
			continue
		}

		if !isZero(field) {
			message += fmt.Sprintf("%s: %v\n", fieldName, field.Interface())
		}
	}

	return message
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Ptr, reflect.Slice:
		return v.IsNil()
	case reflect.Map:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	}
	// Consider other fields non-zero
	return false
}
