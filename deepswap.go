package deepswap

import (
	"reflect"
)

func DeepSwap[T any, V comparable](src T, old V, new V) T {

	srcValue := reflect.ValueOf(src)
	srcType := srcValue.Type()
	oldValue := reflect.ValueOf(old)
	oldType := oldValue.Type()

	oldConverted := oldValue
	if srcValue.Type() != reflect.TypeOf(old) && srcType.Kind() == oldType.Kind() {
		oldConverted = oldValue.Convert(reflect.TypeOf(src))
	}

	switch srcType.Kind() {
	case reflect.Int, reflect.String, reflect.Bool, reflect.Float32, reflect.Float64, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Complex64, reflect.Complex128:

		if srcValue.Interface() == oldConverted.Interface() {
			return reflect.ValueOf(new).Convert(srcType).Interface().(T)
		}
		return src
	case reflect.Pointer:
		if srcValue.IsNil() {
			return src
		}
		newPtrPtr := reflect.New(srcType)
		newPtr := reflect.New(srcType.Elem())
		newValue := DeepSwap(srcValue.Elem().Interface(), old, new)

		newPtr.Elem().Set(reflect.ValueOf(newValue))
		newPtr.Convert(srcType)
		newPtrPtr.Elem().Set(newPtr)
		return newPtrPtr.Elem().Interface().(T)
	case reflect.Interface:
		if srcValue.IsNil() {
			return src
		}
		return reflect.ValueOf(DeepSwap(srcValue.Elem().Interface(), old, new)).Interface().(T)
	case reflect.Map:
		if srcValue.IsNil() {
			return src
		}
		mapCpy := reflect.MakeMapWithSize(srcType, srcValue.Len())
		for _, k := range srcValue.MapKeys() {
			newV := DeepSwap(srcValue.MapIndex(k).Interface(), old, new)
			mapCpy.SetMapIndex(k, reflect.ValueOf(newV))
		}
		return mapCpy.Interface().(T)
	case reflect.Slice:
		if srcValue.IsNil() {
			return src
		}
		sliceCopy := reflect.MakeSlice(srcType, srcValue.Len(), srcValue.Cap())
		for i := 0; i < srcValue.Len(); i++ {
			newV := DeepSwap(srcValue.Index(i).Interface(), old, new)
			sliceCopy.Index(i).Set(reflect.ValueOf(newV))
		}
		return sliceCopy.Interface().(T)
	case reflect.Array:
		if srcValue.Len() == 0 {
			return src
		}
		arrayCopy := reflect.New(srcType).Elem()
		for i := 0; i < srcValue.Len(); i++ {
			newV := DeepSwap(srcValue.Index(i).Interface(), old, new)
			arrayCopy.Index(i).Set(reflect.ValueOf(newV))
		}
		return arrayCopy.Interface().(T)
	case reflect.Struct:
		structCopy := reflect.New(srcType).Elem()
		for i := 0; i < srcValue.NumField(); i++ {
			newFieldVal := DeepSwap(srcValue.Field(i).Interface(), old, new)
			structCopy.FieldByName(srcType.Field(i).Name).Set(reflect.ValueOf(newFieldVal))
		}
		return structCopy.Interface().(T)
	case reflect.UnsafePointer:
		if srcValue.IsNil() {
			return src
		}
		newPtr := reflect.New(srcType.Elem())
		newPtr.Elem().Set(reflect.ValueOf(DeepSwap(reflect.ValueOf(src).Elem().Interface(), old, new)))
		return newPtr.Interface().(T)
	case reflect.Chan, reflect.Func:
		return src
	default:
		return src
	}
}
