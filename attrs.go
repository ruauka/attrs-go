/*
 * Author: ruauka
 *
 * License: MIT (See License file for full text).
 */

// Package attrs_go - tool for working with structure fields. Analog of Python 'getattr' and 'setattr',
// also some useful funcs to change and rounding struct fields.
package attrs_go //nolint:revive,stylecheck

import (
	"cmp"
	"fmt"
	"math"
	"reflect"
)

// float bit size.
const (
	bitSize32 = 32
	bitSize64 = 64
)

// GetAttr - get struct field value.
// 'obj': value param, fields can be ptr or value.
// 'fieldName': value param.
func GetAttr(obj interface{}, fieldName string) (interface{}, error) {
	// to reflect value
	objValue := reflect.ValueOf(obj)
	// struct ptr check
	if objValue.Kind() == reflect.Ptr {
		return nil, ErrPointerStruct
	}
	// is struct check
	if objValue.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}
	// get field value
	field := objValue.FieldByName(fieldName)
	// is field in struct
	if !field.IsValid() {
		return nil, ErrFieldNotInStruct
	}
	// is field exported
	if !field.CanInterface() {
		return nil, ErrUnexportedField
	}
	// field  ptr check
	if field.Kind() == reflect.Ptr {
		return field.Elem().Interface(), nil
	}

	return field.Interface(), nil
}

// SetAttr - set new value on structure field.
// 'obj': ptr struct, fields can be ptr or value.
// 'fieldName', 'newValue': value param.
func SetAttr(obj, newValue interface{}, fieldName string) error {
	var (
		// to reflect value
		objValue = reflect.ValueOf(obj)
		field    reflect.Value
	)
	// struct ptr check
	if objValue.Kind() != reflect.Ptr {
		return ErrNotPointerStruct
	}
	// is struct check
	if objValue.Elem().Kind() != reflect.Struct {
		return ErrNotStruct
	}
	// get field value
	objValue = objValue.Elem().FieldByName(fieldName)
	// field  ptr check
	if objValue.Kind() == reflect.Ptr {
		field = objValue.Elem()
	} else {
		field = objValue
	}
	// is field in struct
	if !field.IsValid() {
		return ErrFieldNotInStruct
	}
	// types check
	if field.Type() != reflect.TypeOf(newValue) {
		return ErrWrongFieldValueType
	}
	// is field exported
	if !field.CanSet() {
		return ErrUnexportedField
	}
	// set value
	field.Set(reflect.ValueOf(newValue))

	return nil
}

// SetStructAttrs - updates current structure fields with the values of the new structure fields.
// 'curObj': ptr struct, fields can be ptr or value.
// 'newObj': value struct, fields can be ptr or value.
func SetStructAttrs(curObj, newObj interface{}) error {
	// to reflect value
	objValue := reflect.ValueOf(newObj)
	// is struct check
	if objValue.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	for i := 0; i < objValue.NumField(); i++ {
		// get field reflect value
		field := objValue.Field(i)
		// is ptr && is ptr nil check
		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}
		// get newObj field name
		fieldName := objValue.Type().Field(i).Name
		// get newObj field value
		fieldValue, err := GetAttr(newObj, fieldName)
		if err != nil {
			return fmt.Errorf("err in GetAttr: %w", err)
		}
		// get curObj field value
		if err := SetAttr(curObj, fieldValue, fieldName); err != nil {
			return fmt.Errorf("err in SetAttr: %w", err)
		}
	}

	return nil
}

// Round - float64 | float32 rounder to certain precision.
func Round[V float64 | float32](value V, precision int) V {
	return V(math.Round(float64(value)*(math.Pow10(precision))) / math.Pow10(precision))
}

// iterRound - round any float. Private func.
func iterRound(field reflect.Value, precision int, bitSize int) {
	for j := 0; j < field.Len(); j++ {
		field := field.Index(j)
		// bitSize check
		if bitSize == bitSize64 {
			field.Set(reflect.ValueOf(Round(field.Float(), precision)))
			continue
		}

		field.Set(reflect.ValueOf(float32(Round(field.Float(), precision))))
	}
}

// RoundFloatStruct - round up float struct fields to certain precision
// Constraint: simple floats, array and slice.
// 'obj': ptr struct, fields can be value, not ptr.
// 'precision': round to.
func RoundFloatStruct(obj interface{}, precision int) error {
	objValue := reflect.ValueOf(obj)
	// struct ptr check
	if objValue.Kind() != reflect.Ptr {
		return ErrNotPointerStruct
	}
	// is struct check
	if objValue.Elem().Kind() != reflect.Struct {
		return ErrNotStruct
	}
	// get value from ptr
	objValue = objValue.Elem()

	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Field(i)
		// is field exported
		if !field.CanSet() {
			return ErrUnexportedField
		}
		// float types check
		switch field.Kind() {
		// simple float
		case reflect.Float64, reflect.Float32:
			if field.Kind() == reflect.Float64 {
				field.Set(reflect.ValueOf(Round(field.Float(), precision)))
				break
			}
			field.Set(reflect.ValueOf(float32(Round(field.Float(), precision))))
		// array and slice float
		case reflect.Array, reflect.Slice:
			if field.Len() == 0 {
				break
			}
			if field.Index(0).Kind() == reflect.Float64 {
				iterRound(field, precision, bitSize64)
				break
			}
			iterRound(field, precision, bitSize32)
		}
	}

	return nil
}

// Intersection - intersection of two arrays. Returns new slice.
func Intersection[T cmp.Ordered](left, right []T) []T {
	var (
		minimum = min(len(left), len(right))
		out     = make([]T, 0, minimum)
		check   = make(map[T]struct{}, minimum)
	)

	for _, i := range left {
		for _, j := range right {
			if i == j && check[i] == struct{}{} {
				out = append(out, i)
				check[i] = struct{}{}
			}
		}
	}

	return out
}

// SlicesConcat - concatenation of multiple slices.
func SlicesConcat[T any](slices ...[]T) []T {
	var length, idx int

	for _, slice := range slices {
		length += len(slice)
	}

	res := make([]T, length)

	for _, s := range slices {
		idx += copy(res[idx:], s)
	}

	return res
}
