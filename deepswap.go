package deepswap

import "reflect"

func DeepSwap[T any, V comparable](src T, old V, new V) any {

	srcValue := reflect.ValueOf(src)
	srcType := srcValue.Type()

	switch srcType.Kind() {
	case reflect.Int, reflect.String, reflect.Bool, reflect.Float32, reflect.Float64, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Complex64, reflect.Complex128:
		if srcType == reflect.TypeOf(old) {
			typecastSrc := any(old).(V)
			if typecastSrc == old {
				return any(new)
			}
		}
		return src
	case reflect.Pointer:
		if srcValue.IsNil() {
			return src // Return nil pointer as is
		}
		// Recursively modify the dereferenced value
		newPtr := reflect.New(srcType.Elem()) // Create a new pointer of the same type
		newPtr.Elem().Set(reflect.ValueOf(DeepSwap(srcValue.Elem().Interface(), old, new)))
		return newPtr.Interface()
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
			// newK := reflect.ValueOf(DeepSwap(k.Interface(), old, new))
			newV := DeepSwap(srcValue.MapIndex(k).Interface(), old, new)
			mapCpy.SetMapIndex(k, reflect.ValueOf(newV))
		}
		return mapCpy.Interface()
	case reflect.Slice:
		if srcValue.IsNil() {
			return src
		}
		sliceCopy := reflect.MakeSlice(srcType, srcValue.Len(), srcValue.Cap())
		for i := 0; i < srcValue.Len(); i++ {
			newV := DeepSwap(srcValue.Index(i).Interface(), old, new)
			sliceCopy.Index(i).Set(reflect.ValueOf(newV))
		}
		return sliceCopy.Interface()
	case reflect.Array:
		if srcValue.Len() == 0 {
			return src
		}
		arrayCopy := reflect.New(srcType).Elem()
		for i := 0; i < srcValue.Len(); i++ {
			newV := DeepSwap(srcValue.Index(i).Interface(), old, new)
			arrayCopy.Index(i).Set(reflect.ValueOf(newV))
		}
		return arrayCopy.Interface()
	case reflect.Struct:
		structCopy := reflect.New(reflect.StructOf(reflect.VisibleFields(srcType))).Elem()
		for i := 0; i < srcValue.NumField(); i++ {
			newFieldVal := DeepSwap(srcValue.Field(i).Interface(), old, new)
			structCopy.FieldByName(srcType.Field(i).Name).Set(reflect.ValueOf(newFieldVal))
		}
		return structCopy.Interface()
	case reflect.Chan:
		return src
	case reflect.Func:
		return src
	case reflect.UnsafePointer:
		if srcValue.IsNil() {
			return src
		}
		newPtr := reflect.New(srcType.Elem())
		newPtr.Elem().Set(reflect.ValueOf(DeepSwap(reflect.ValueOf(src).Elem().Interface(), old, new)))
		return newPtr.Interface()
	default:
		// Direct replacement if it matches `old`
		panic("type not known")
	}
}
