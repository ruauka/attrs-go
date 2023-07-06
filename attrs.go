package attrs_go

import (
	"errors"
	"log"
	"reflect"
)

var (
	ErrNotStruct           = errors.New("not a struct")
	ErrFieldNotInStruct    = errors.New("field not in struct")
	ErrUnexportedField     = errors.New("field not exported")
	ErrWrongFieldValueType = errors.New("wrong field value type")
	ErrNotPointer          = errors.New("struct passed not by pointer")
)

// setValue - set struct new field.
func setValue(obj interface{}, fieldName string, newValue interface{}) error {
	objValue := reflect.ValueOf(obj)

	if objValue.Kind() != reflect.Ptr {
		return ErrNotPointer
	}

	if objValue.Elem().Kind() != reflect.Struct {
		return ErrNotStruct
	}

	field := objValue.Elem().FieldByName(fieldName)
	if !field.IsValid() {
		return ErrFieldNotInStruct
	}

	if field.Type() != reflect.TypeOf(newValue) {
		return ErrWrongFieldValueType
	}

	if !field.CanSet() {
		return ErrUnexportedField
	}

	field.Set(reflect.ValueOf(newValue))

	return nil
}

// getValue - get struct field value.
func getValue(obj interface{}, fieldName string) (interface{}, error) {
	objValue := reflect.ValueOf(obj)

	if objValue.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	field := reflect.ValueOf(obj).FieldByName(fieldName)
	if !field.IsValid() {
		return nil, ErrFieldNotInStruct
	}

	return field.Interface(), nil
}

func SetStructAttrs(curObj, newObj interface{}) {
	elem := reflect.ValueOf(newObj)

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if field.IsNil() {
			continue
		}

		fieldName := elem.Type().Field(i).Name

		fieldValue, err := getValue(newObj, fieldName)
		if err != nil {
			log.Fatalf("Err in GetValue: %v", err)
		}

		if err := setValue(curObj, fieldName, fieldValue); err != nil {
			log.Fatalf("Err in SetValue: %v", err)
		}
	}
}
