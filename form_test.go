package form

import (
	"fmt"
	"testing"
	//"net/http"
	"errors"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
type Human struct {
	Name string `form:"name"`
	Age  int    `form:"age"`
}

func (this Human) NameValidator(n string) []error {
	return []error{errors.New("Name 字段错误1."), errors.New("Name 字段错误2.")}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
type Class struct {
	ClassName string `form:"class_name"`
}

func (this Class) ClassNameValidator(n string) error {
	return errors.New("ClassName 字段错误.")
}

////////////////////////////////////////////////////////////////////////////////////////////////////
type Student struct {
	*Human
	Number int `form:"number"`
	Class  Class
}

func (this Student) NumberValidator(n int) error {
	return errors.New("Number 字段错误.")
}

////////////////////////////////////////////////////////////////////////////////////////////////////
var formData = map[string][]string{"name": []string{"Yangfeng"}, "age": []string{"12"}, "number": []string{"9"}, "class_name":[]string{"class one"}}

func TestBindPoint(t *testing.T) {
	fmt.Println("===== bind pointer =====")
	var s *Student
	Bind(&s, formData)

	fmt.Println("数据验证错误:", Validate(s))

	if s != nil {
		fmt.Println(s.Name, s.Age, s.Number, s.Class.ClassName)
	}
}

func TestBindStruct(t *testing.T) {
	fmt.Println("===== bind struct =====")
	var s Student
	//Bind(&s, formData)

	fmt.Println("数据验证错误:", Validate(s))

	//fmt.Println(s.Name, s.Age, s.Number, s.Class.ClassName)
}

type People struct {
	Form
	Name      string  `form:"name"`
	Age       int     `form:"age"`
	Undefined string  `form:"undefined"` // 表单中没有的字段，其不会出现在 CleanData 中
}

//func TestCleanData(t *testing.T) {
//	fmt.Println("===== bind with clean data =====")
//	var p People
//	Bind(&p, formData)
//	fmt.Println(p.Name, p.Age, p.CleanData)
//}
//
//func TestRequest(t *testing.T) {
//	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
//		var p People
//		BindWithRequest(req, &p)
//		writer.Write([]byte(fmt.Sprintf("name: %s  age: %d", p.Name, p.Age)))
//	})
//	http.ListenAndServe(":8000", nil)
//}