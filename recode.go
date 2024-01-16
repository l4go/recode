package recode

import (
	"reflect"
)

type AnyRebuilder[T any] interface {
	RebuildByType(param T) error
}

func deepIndirect(rv reflect.Value) reflect.Value {
	for rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}

	return rv
}

func RecursiveRebuild[T any](v any, param T) error {
	return recursiveRebuild[T](reflect.ValueOf(v), param)
}

func rebuild[T any](rv reflect.Value, param T) error {
	if !rv.CanAddr() {
		return nil
	}
	rva := rv.Addr()
	if !rva.CanInterface() {
		return nil
	}

	rvi := rva.Interface()
	v, ok := rvi.(AnyRebuilder[T])
	if ok {
		err := v.RebuildByType(param)
		if err != nil {
			return err
		}
	}

	return nil
}

func recursiveRebuild[T any](rv reflect.Value, param T) error {
	rv = deepIndirect(rv)
	err := rebuild[T](rv, param)
	if err != nil {
		return err
	}

	switch rv.Kind() {
	case reflect.Struct:
		typ := rv.Type()
		for i := 0; i < typ.NumField(); i++ {
			err := recursiveRebuild[T](rv.Field(i), param)
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		it := rv.MapRange()
		for it.Next() {
			recursiveRebuild[T](it.Value(), param)
		}
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			err := recursiveRebuild[T](rv.Index(i), param)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
