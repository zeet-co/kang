//thanks chatGPT
package controller

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func setField(obj interface{}, name, value string) error {
	structValue := reflect.ValueOf(obj).Elem()
	fieldValue := structValue.FieldByName(name)

	if !fieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	fieldType := fieldValue.Type()
	isPtr := fieldType.Kind() == reflect.Ptr
	var actualType reflect.Type
	if isPtr {
		actualType = fieldType.Elem()
	} else {
		actualType = fieldType
	}

	switch actualType.Kind() {
	case reflect.String:
		if isPtr {
			fieldValue.Set(reflect.ValueOf(&value))
		} else {
			fieldValue.SetString(value)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, actualType.Bits())
		if err != nil {
			return err
		}
		if isPtr {
			intPtr := reflect.New(actualType)
			intPtr.Elem().SetInt(intValue)
			fieldValue.Set(intPtr)
		} else {
			fieldValue.SetInt(intValue)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, actualType.Bits())
		if err != nil {
			return err
		}
		if isPtr {
			uintPtr := reflect.New(actualType)
			uintPtr.Elem().SetUint(uintValue)
			fieldValue.Set(uintPtr)
		} else {
			fieldValue.SetUint(uintValue)
		}
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, actualType.Bits())
		if err != nil {
			return err
		}
		if isPtr {
			floatPtr := reflect.New(actualType)
			floatPtr.Elem().SetFloat(floatValue)
			fieldValue.Set(floatPtr)
		} else {
			fieldValue.SetFloat(floatValue)
		}
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		if isPtr {
			fieldValue.Set(reflect.ValueOf(&boolValue))
		} else {
			fieldValue.SetBool(boolValue)
		}
	default:
		return fmt.Errorf("Unsupported field type %s", fieldType)
	}
	return nil
}

// func setField(obj interface{}, name string, value string) error {
// structValue := reflect.ValueOf(obj).Elem()
// fieldValue := structValue.FieldByName(name)

// if !fieldValue.IsValid() {
// return fmt.Errorf("No such field: %s in obj", name)
// }

// if !fieldValue.CanSet() {
// return fmt.Errorf("Cannot set %s field value", name)
// }

// fieldType := fieldValue.Type()
// switch fieldType.Kind() {
// case reflect.String:
// fieldValue.SetString(value)
// case reflect.Int:
// intValue, err := strconv.Atoi(value)
// if err != nil {
// return err
// }
// fieldValue.SetInt(int64(intValue))
// // Add more cases here for other types as needed
// default:
// return fmt.Errorf("Unsupported field type %s", fieldType)
// }
// return nil
// }

func getJSONTagToFieldNameMap(obj interface{}) map[string]string {
	tagToName := make(map[string]string)
	val := reflect.ValueOf(obj).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			tagToName[jsonTag] = field.Name
		}
	}
	return tagToName
}

func assignValues(obj interface{}, values map[string]string) (error, bool) {
	tagToName := getJSONTagToFieldNameMap(obj)
	anyFieldSet := false
	var err error

	for key, value := range values {
		keys := strings.Split(key, ".")

		if len(keys) > 1 {
			// Nested struct
			nestedFieldName, ok := tagToName[keys[0]]
			if !ok {
				continue
				// return fmt.Errorf("No such json tag: %s in obj", keys[0])
			}

			nestedField := reflect.ValueOf(obj).Elem().FieldByName(nestedFieldName)
			if !nestedField.IsValid() {
				return fmt.Errorf("Invalid field for json tag: %s", keys[0]), anyFieldSet
			}

			// Check if the nested field is a pointer to a struct
			if nestedField.Kind() == reflect.Ptr && nestedField.Type().Elem().Kind() == reflect.Struct {
				if nestedField.IsNil() {
					// If the struct pointer is nil, create a new instance
					newStruct := reflect.New(nestedField.Type().Elem())
					nestedField.Set(newStruct)
				}
				// Recursively assign values to the nested struct
				err, anyFieldSet = assignValues(nestedField.Interface(), map[string]string{strings.Join(keys[1:], "."): value})
				if err != nil {
					return err, anyFieldSet
				}
			} else if nestedField.Kind() == reflect.Struct {
				// Handle non-pointer nested structs
				err, anyFieldSet = assignValues(nestedField.Addr().Interface(), map[string]string{strings.Join(keys[1:], "."): value})
				if err != nil {
					return err, anyFieldSet
				}
			} else {
				return fmt.Errorf("Field %s is not a struct or pointer to a struct", nestedFieldName), anyFieldSet
			}
		} else {
			fieldName, ok := tagToName[key]
			if !ok {
				continue
				// return fmt.Errorf("No such json tag: %s in obj", key)
			}
			err := setField(obj, fieldName, value)
			if err != nil {
				return err, anyFieldSet
			}
			anyFieldSet = true
		}
	}
	return nil, anyFieldSet
}
