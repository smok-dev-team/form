package form

import (
	"reflect"
	"strconv"
	"fmt"
)

func Bind(obj interface{}, form map[string][]string) (err error) {
	fmt.Println("in")
	var objValue = reflect.ValueOf(obj)
	var objType = reflect.TypeOf(obj)

	var objValueKind = objValue.Kind()
	for {
		if objValueKind == reflect.Ptr && objValue.IsNil() {
			objValue.Set(reflect.New(objType.Elem()))
		}

		if objValueKind == reflect.Ptr {
			objValue = objValue.Elem()
			objType = objType.Elem()
			objValueKind = objValue.Kind()
			continue
		}
		break
	}

	var cleanDataValue = objValue.FieldByName("CleanData")
	if cleanDataValue.IsValid() && cleanDataValue.IsNil() {
		cleanDataValue.Set(reflect.MakeMap(cleanDataValue.Type()))
	}
	return mapForm(objType, objValue, cleanDataValue, form)
}

func mapForm(objType reflect.Type, objValue, cleanDataValue reflect.Value, form map[string][]string) (err error) {
	var numField = objType.NumField()
	for i:=0; i< numField; i++ {
		var fieldType = objType.Field(i)
		var fieldValue = objValue.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		var tag = fieldType.Tag.Get("form")

		if tag == "" {
			tag = fieldType.Name

			if fieldValue.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}
				fieldValue = fieldValue.Elem()
			}

			if fieldValue.Kind() == reflect.Struct {
				err = mapForm(fieldValue.Addr().Type().Elem(), fieldValue, cleanDataValue, form)
				if err != nil {
					return err
				}
				continue
			}
		} else if tag == "-" {
			continue
		}

		var values, exists = form[tag]
		if !exists {
			continue
		}

		err = setValueForField(fieldType.Type.Kind(), values[0], fieldValue)
		if err != nil {
			return err
		}
		if cleanDataValue.IsValid() {
			cleanDataValue.SetMapIndex(reflect.ValueOf(tag), fieldValue)
		}
	}
	return err
}

func setValueForField(fieldTypeKind reflect.Kind, v string, fieldValue reflect.Value) (error) {
	switch fieldTypeKind {
	case reflect.Int:
		return setIntField(v, 0, fieldValue)
	case reflect.Int8:
		return setIntField(v, 8, fieldValue)
	case reflect.Int16:
		return setIntField(v, 16, fieldValue)
	case reflect.Int32:
		return setIntField(v, 32, fieldValue)
	case reflect.Int64:
		return setIntField(v, 64, fieldValue)
	case reflect.Uint:
		return setUintField(v, 0, fieldValue)
	case reflect.Uint8:
		return setUintField(v, 8, fieldValue)
	case reflect.Uint16:
		return setUintField(v, 16, fieldValue)
	case reflect.Uint32:
		return setUintField(v, 32, fieldValue)
	case reflect.Uint64:
		return setUintField(v, 64, fieldValue)
	case reflect.Float32:
		return setFloatField(v, 32, fieldValue)
	case reflect.Float64:
		return setFloatField(v, 64, fieldValue)
	case reflect.Bool:
		return setBoolField(v, 0, fieldValue)
	case reflect.String:
		fieldValue.SetString(v)
	}

	return nil
}

func setIntField(v string, bitSize int, vf reflect.Value) (error) {
	if v == "" {
		v = "0"
	}
	var iv, err = strconv.ParseInt(v, 10, bitSize)
	if err == nil {
		vf.SetInt(iv)
	}
	return err
}

func setUintField(v string, bitSize int, vf reflect.Value) (error) {
	if v == "" {
		v = "0"
	}
	var iv, err = strconv.ParseUint(v, 10, bitSize)
	if err == nil {
		vf.SetUint(iv)
	}
	return err
}

func setBoolField(v string, bitSize int, vf reflect.Value) (error) {
	if v == "" {
		v = "false"
	}
	var iv, err = strconv.ParseBool(v)
	if err == nil {
		vf.SetBool(iv)
	}
	return err
}

func setFloatField(v string, bitSize int, vf reflect.Value) (error) {
	if v == "" {
		v = "0.0"
	}
	var iv, err = strconv.ParseFloat(v, bitSize)
	if err == nil {
		vf.SetFloat(iv)
	}
	return err
}