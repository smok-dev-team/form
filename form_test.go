package form

import (
	"fmt"
	"net/http"
	"testing"
)

type Human struct {
	Name string `form:"name"`
	Age  int    `form:"age"`
}

type Class struct {
	ClassName string `form:"class_name"`
}

type Student struct {
	Human
	Number int `form:"number"`
	Class  Class
}

var formData = map[string][]string{"name": []string{"Yangfeng"}, "age": []string{"120"}, "number": []string{"9"}, "class_name": []string{"class one"}}

func TestBindPoint(t *testing.T) {
	fmt.Println("===== bind pointer =====")
	var s *Student
	Bind(formData, &s)

	if s != nil {
		fmt.Println(s.Name, s.Age, s.Number, s.Class.ClassName)
	}
}

func TestBindStruct(t *testing.T) {
	fmt.Println("===== bind struct =====")
	var s Student
	Bind(formData, &s)

	fmt.Println(s.Name, s.Age, s.Number, s.Class.ClassName)
}

type People struct {
	Form
	Name      string `form:"name"`
	Age       int    `form:"age"`
	Gender    int
	Undefined string `form:"undefined"` // 表单中没有的字段，其不会出现在 CleanData 中
}

// 可以为属性添加一个提供默认值的方法
func (this People) DefaultGender() int {
	return 1
}

func TestCleanData(t *testing.T) {
	fmt.Println("===== bind with clean data =====")
	var p People
	Bind(formData, &p)
	fmt.Println(p.Name, p.Age, p.Gender, p.CleanedData)
}

func TestRequest(t *testing.T) {
	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		var p People
		BindWithRequest(req, &p)
		writer.Write([]byte(fmt.Sprintf("name: %s  age: %d", p.Name, p.Age)))
	})
	http.ListenAndServe(":8000", nil)
}
