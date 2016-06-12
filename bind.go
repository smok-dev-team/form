package form

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	k_FORM_TAG                 = "form"
	k_FORM_NO_TAG              = "-"
	k_FORM_CLEANED_DATA        = "CleanedData"
	k_FORM_DEFAULT_FUNC_PREFIX = "Default"
)

func BindWithRequest(request *http.Request, result interface{}) (err error) {
	err = request.ParseForm()
	if err != nil {
		return err
	}

	var contentType = request.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		request.ParseMultipartForm(32 << 20)
	}

	err = Bind(request.Form, result)
	return err
}

// var err = Bind(&result, data)
func Bind(source map[string][]string, result interface{}) (err error) {
	var objValue = reflect.ValueOf(result)
	var objType = reflect.TypeOf(result)
	var objValueKind = objValue.Kind()

	if objValueKind == reflect.Struct {
		return errors.New("obj is struct")
	}
	if objValue.IsNil() {
		return errors.New("obj is nil")
	}

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

	var cleanDataValue = objValue.FieldByName(k_FORM_CLEANED_DATA)
	if cleanDataValue.IsValid() && cleanDataValue.IsNil() {
		cleanDataValue.Set(reflect.MakeMap(cleanDataValue.Type()))
	}
	return mapForm(objType, objValue, cleanDataValue, source)
}

func mapForm(objType reflect.Type, objValue, cleanDataValue reflect.Value, source map[string][]string) error {
	var numField = objType.NumField()
	for i := 0; i < numField; i++ {
		var fieldStruct = objType.Field(i)
		var fieldValue = objValue.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		var tag = fieldStruct.Tag.Get(k_FORM_TAG)

		if tag == "" {
			tag = fieldStruct.Name

			if fieldValue.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}
				fieldValue = fieldValue.Elem()
			}

			if fieldValue.Kind() == reflect.Struct {
				if err := mapForm(fieldValue.Addr().Type().Elem(), fieldValue, cleanDataValue, source); err != nil {
					return err
				}
				continue
			}
		} else if tag == k_FORM_NO_TAG {
			continue
		}

		var values, exists = source[tag]
		if !exists {
			var mName = k_FORM_DEFAULT_FUNC_PREFIX + fieldStruct.Name
			var mValue = objValue.MethodByName(mName)
			if mValue.IsValid() == false {
				if objValue.CanAddr() {
					mValue = objValue.Addr().MethodByName(mName)
				}
			}

			if mValue.IsValid() {
				var rList = mValue.Call(nil)
				if cleanDataValue.IsValid() {
					cleanDataValue.SetMapIndex(reflect.ValueOf(tag), rList[0])
				}
				fieldValue.Set(rList[0])
			}
			continue
		}

		if fieldValue.Kind() == reflect.Slice {
			var valueLen = len(values)
			var sKind = fieldValue.Type().Elem().Kind()
			var s = reflect.MakeSlice(fieldStruct.Type, valueLen, valueLen)
			for i := 0; i < valueLen; i++ {
				if err := setValueForField(sKind, values[i], s.Index(i)); err != nil {
					return err
				}
			}
			objValue.Field(i).Set(s)
		} else {
			if err := setValueForField(fieldStruct.Type.Kind(), values[0], fieldValue); err != nil {
				return err
			}
		}
		if cleanDataValue.IsValid() {
			cleanDataValue.SetMapIndex(reflect.ValueOf(tag), fieldValue)
		}
	}
	return nil
}

func setValueForField(fieldTypeKind reflect.Kind, v string, fieldValue reflect.Value) error {
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

func setIntField(v string, bitSize int, vf reflect.Value) error {
	if v == "" {
		v = "0"
	}
	var iv, err = strconv.ParseInt(v, 10, bitSize)
	if err == nil {
		vf.SetInt(iv)
	}
	return err
}

func setUintField(v string, bitSize int, vf reflect.Value) error {
	if v == "" {
		v = "0"
	}
	var iv, err = strconv.ParseUint(v, 10, bitSize)
	if err == nil {
		vf.SetUint(iv)
	}
	return err
}

func setBoolField(v string, bitSize int, vf reflect.Value) error {
	if v == "" {
		v = "false"
	}
	var iv, err = strconv.ParseBool(v)
	if err == nil {
		vf.SetBool(iv)
	}
	return err
}

func setFloatField(v string, bitSize int, vf reflect.Value) error {
	if v == "" {
		v = "0.0"
	}
	var iv, err = strconv.ParseFloat(v, bitSize)
	if err == nil {
		vf.SetFloat(iv)
	}
	return err
}
