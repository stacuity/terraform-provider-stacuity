// Copyright (c) HashiCorp, Inc.

package stacuity

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func GetReflectValues(src interface{}, destPtr interface{}) (srcType reflect.Type, srcVal reflect.Value, destVal reflect.Value) {
	if reflect.TypeOf(src).Kind() == reflect.Ptr {
		srcVal = reflect.ValueOf(src).Elem()
		srcType = reflect.TypeOf(src).Elem()
	} else {
		srcVal = reflect.ValueOf(src)
		srcType = reflect.TypeOf(src)
	}

	if reflect.TypeOf(destPtr).Kind() == reflect.Ptr {
		destVal = reflect.ValueOf(destPtr).Elem()
	} else {
		destVal = reflect.ValueOf(destPtr)
	}

	return
}

// JC 05/11/24 - Created to speed up mapping to and from terraform to/from our API.
func ConvertFromAPI(src interface{}, destPtr interface{}) error {
	// Check if src is nil pointer and return early if true
	if reflect.ValueOf(src).Kind() == reflect.Ptr && reflect.ValueOf(src).IsNil() {
		return nil
	}

	srcType, srcVal, destVal := GetReflectValues(src, destPtr)

	if destVal.Kind() == reflect.Slice {
		if srcVal.Kind() != reflect.Slice {
			return errors.New("source and destination must be a slice")
		}

		destSlice := reflect.MakeSlice(destVal.Type(), srcVal.Len(), srcVal.Cap())
		for i := 0; i < srcVal.Len(); i++ {
			srcElem := srcVal.Index(i)
			destElem := reflect.New(destVal.Type().Elem()).Elem()
			if err := ConvertFromAPI(srcElem.Interface(), destElem.Addr().Interface()); err != nil {
				return err
			}
			destSlice.Index(i).Set(destElem)
		}
		destVal.Set(destSlice)
		return nil
	}

	if destVal.Kind() != reflect.Struct {
		return errors.New("destPtr " + destVal.String() + " destination should be a struct not " + destVal.Kind().String())
	}

	if srcType.Kind() == reflect.String {
		if destVal.Kind() == reflect.Ptr && destVal.IsNil() {
			newDest := reflect.New(destVal.Type().Elem())
			destVal.Set(newDest)
		}

		var sv types.String
		if srcType.Kind() == reflect.String {
			sv = types.StringValue(srcVal.String())
		} else {
			sv = types.StringPointerValue(src.(*string))
		}

		destVal.Set(reflect.ValueOf(sv))
		return nil
	}

	if srcType.Kind() == reflect.Int32 {
		if destVal.Kind() == reflect.Ptr && destVal.IsNil() {
			newDest := reflect.New(destVal.Type().Elem())
			destVal.Set(newDest)
		}
		sv := types.Int32PointerValue(src.(*int32))
		destVal.Set(reflect.ValueOf(sv))
		return nil
	}

	for i := 0; i < srcType.NumField(); i++ {
		fieldType := srcType.Field(i)
		fieldName := fieldType.Name
		srcField := srcVal.Field(i)
		destField := destVal.FieldByName(fieldName)
		srcData := srcVal.Field(i).Interface()

		if srcData == nil || !srcField.IsValid() || !destField.IsValid() || (srcField.Kind() == reflect.Ptr && srcField.IsNil()) {
			continue // skip if no match found
		}

		destType := destField.Type()

		switch reflect.TypeOf(srcData).Kind() {
		case reflect.String:
			sv := types.StringValue(srcData.(string))
			destField.Set(reflect.ValueOf(sv))
		case reflect.Int32:
			sv := types.Int32Value(srcData.(int32))
			destField.Set(reflect.ValueOf(sv))
		case reflect.Bool:
			sv := types.BoolValue(srcData.(bool))
			destField.Set(reflect.ValueOf(sv))
		case reflect.Slice:
			destSlice := reflect.MakeSlice(destType, srcField.Len(), srcField.Cap())
			for i := 0; i < srcField.Len(); i++ {
				srcElem := srcField.Index(i)
				destElem := reflect.New(destType.Elem()).Elem()
				err := ConvertFromAPI(srcElem.Interface(), destElem.Addr().Interface())
				if err != nil {
					return err
				}
				destSlice.Index(i).Set(destElem)
			}
			destField.Set(destSlice)
		case reflect.Struct:
			if destField.Kind() == reflect.Ptr && destField.IsNil() {
				newDest := reflect.New(destField.Type().Elem())
				destField.Set(newDest)
				err := ConvertFromAPI(srcData, newDest.Interface())
				if err != nil {
					return err
				}
			}

			if destField.Kind() == reflect.Ptr {
				newDestData := reflect.New(destType.Elem()).Interface()
				err := ConvertFromAPI(srcData, newDestData)
				if err != nil {
					return err
				}
				destField.Set(reflect.ValueOf(newDestData))
			}

			if destField.Kind() == reflect.Struct {
				err := ConvertFromAPI(srcData, destField.Addr().Interface())
				if err != nil {
					return err
				}
			}
		case reflect.Ptr:
			if destType.Kind() == reflect.Slice {
				srcField = srcField.Elem()
				destSlice := reflect.MakeSlice(destType, srcField.Len(), srcField.Cap())
				for i := 0; i < srcField.Len(); i++ {
					srcElem := srcField.Index(i)
					destElem := reflect.New(destType.Elem().Elem())
					err := ConvertFromAPI(srcElem.Interface(), destElem.Interface())
					if err != nil {
						return err
					}
					destSlice.Index(i).Set(destElem)
				}
				destField.Set(destSlice)
			} else if destType.Kind() == reflect.Struct {
				err := ConvertFromAPI(srcData, destField.Addr().Interface())
				if err != nil {
					return err
				}
			} else {
				newDestData := reflect.New(destType.Elem()).Interface()
				err := ConvertFromAPI(srcData, newDestData)

				if err != nil {
					return err
				}
				destField.Set(reflect.ValueOf(newDestData))
			}

		default:
			return fmt.Errorf("unsupported type for field %s: %T", fieldName, srcData)
		}
	}
	return nil
}

func ConvertToAPI(src, dest interface{}) error {
	// Doesn't need reflecting
	if reflect.ValueOf(src).Kind() == reflect.Ptr && reflect.ValueOf(src).IsNil() {
		return nil
	}

	srcType, srcVal, destVal := GetReflectValues(src, dest)

	if destVal.Kind() == reflect.Slice {
		if srcVal.Kind() != reflect.Slice {
			return errors.New("source should be a slice when destination is a slice")
		}

		destSlice := reflect.MakeSlice(destVal.Type(), srcVal.Len(), srcVal.Cap())
		for i := 0; i < srcVal.Len(); i++ {
			srcElem := srcVal.Index(i)
			destElem := reflect.New(destVal.Type().Elem()).Elem()
			if err := ConvertToAPI(srcElem.Interface(), destElem.Addr().Interface()); err != nil {
				return err
			}
			destSlice.Index(i).Set(destElem)
		}
		destVal.Set(destSlice)
		return nil
	}

	if destVal.Kind() != reflect.Struct {
		return errors.New("destination should be a struct not " + destVal.Kind().String())
	}

	for i := 0; i < srcType.NumField(); i++ {
		fieldType := srcType.Field(i)
		fieldName := fieldType.Name
		srcField := srcVal.Field(i)
		destField := destVal.FieldByName(fieldName)

		if !destField.IsValid() {
			continue
		}

		destType := destField.Type()

		srcData := srcVal.Field(i).Interface()
		switch reflect.TypeOf(srcData).Kind() {
		case reflect.Slice:

			destSlice := reflect.MakeSlice(destType, srcField.Len(), srcField.Cap())

			//Support basic string slices
			if srcField.Type().Elem().String() == "basetypes.StringValue" && destType.Elem().Kind() == reflect.String {
				for i := 0; i < srcField.Len(); i++ {
					srcElem := srcField.Index(i).Interface().(basetypes.StringValue)
					destSlice.Index(i).SetString(srcElem.ValueString())
				}
			} else {

				for i := 0; i < srcField.Len(); i++ {
					srcElem := srcField.Index(i)
					destElem := reflect.New(destType.Elem()).Elem()

					if reflect.TypeOf(srcElem).Kind() == reflect.Pointer || reflect.TypeOf(srcElem).Kind() == reflect.Struct {
						err := ConvertToAPI(srcElem.Interface(), destElem.Addr().Interface())
						if err != nil {
							return err
						}
					}

					destSlice.Index(i).Set(destElem)
				}
			}
			destField.Set(destSlice)
		case reflect.Struct:
			var tfType = reflect.TypeOf(srcData).String()
			switch tfType {
			case "basetypes.StringValue":
				if destType.Kind() == reflect.String {
					sv := srcData.(basetypes.StringValue)
					destField.Set(reflect.ValueOf(sv.ValueString()))
				} else if destType.Kind() == reflect.Ptr {
					sv := srcData.(basetypes.StringValue)
					destField.Set(reflect.ValueOf(sv.ValueStringPointer()))
				}
			case "basetypes.Int64Value":
				sv := srcData.(basetypes.Int64Value)
				destField.Set(reflect.ValueOf(sv.ValueInt64()))
			case "basetypes.Int32Value":
				if destType.Kind() == reflect.Int32 {
					sv := srcData.(basetypes.Int32Value)
					destField.Set(reflect.ValueOf(sv.ValueInt32()))
				} else if destType.Kind() == reflect.Ptr {
					sv := srcData.(basetypes.Int32Value)
					destField.Set(reflect.ValueOf(sv.ValueInt32Pointer()))
				}
			case "basetypes.BoolValue":
				sv := srcData.(basetypes.BoolValue)
				destField.Set(reflect.ValueOf(sv.ValueBool()))
			default:
				if err := ConvertToAPI(srcField.Interface(), destField.Addr().Interface()); err != nil {
					return err
				}
			}
		case reflect.Ptr:
			newDestData := reflect.New(destType.Elem()).Interface()
			if err := ConvertToAPI(srcData, newDestData); err != nil {
				return err
			}
			destField.Set(reflect.ValueOf(newDestData))
		default:
			return fmt.Errorf("unsupported type for field %s: got %s could be %s", fieldName, reflect.TypeOf(srcData).Kind(), reflect.TypeOf(srcData))
		}
	}

	return nil
}
