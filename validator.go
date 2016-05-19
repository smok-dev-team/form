package form

import (
	"reflect"
	"fmt"
)

const (
	K_VALIDATOR_FUNC_SUFFIX = "Validator"
)

type ValidatorError struct {
	Field   string
	Code    int
	Message string
}

func (this *ValidatorError) Error() string {
	return fmt.Sprintf("[%s]%d:%s", this.Field, this.Code, this.Message)
}

func Validate(obj interface{}) (err []error) {
	var objType = reflect.TypeOf(obj)
	var objValue = reflect.ValueOf(obj)

	for {
		if objValue.Kind() == reflect.Ptr {
			objValue = objValue.Elem()
			objType = objType.Elem()
			continue
		}
		break
	}

	if !objValue.IsValid() {
		return nil
	}

	var errMap = make(map[string][]error)
	validate(objType, objValue, errMap)

	var errList []error
	if len(errMap) > 0 {
		errList = make([]error, 0, 0)
		for _, value := range errMap {
			errList = append(errList, value...)
		}
	}
	return errList
}

func validate(objType reflect.Type, objValue reflect.Value, errMap map[string][]error) {
	var numField = objType.NumField()
	for i:=0; i<numField; i++ {
		var fieldStruct = objType.Field(i)
		var fieldValue = objValue.Field(i)

		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		if fieldValue.Kind() == reflect.Struct {
			validate(fieldValue.Type(), fieldValue, errMap)
			continue
		}

		var mName  = fieldStruct.Name + K_VALIDATOR_FUNC_SUFFIX
		var mValue = objValue.MethodByName(mName)
		if mValue.IsValid() {
			var eList = mValue.Call([]reflect.Value{fieldValue})

			if !eList[0].IsNil() {
				if eList[0].Kind() == reflect.Slice {
					errMap[fieldStruct.Name] = eList[0].Interface().([]error)
				} else {
					errMap[fieldStruct.Name] = []error{eList[0].Interface().(error)}
				}
			}
		}
	}
}