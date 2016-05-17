package form

import (
	"fmt"
	"testing"
)

type Human struct {
	*Person
	CleanData map[string]interface{}
	Name string `bson:"name" form:"name"`
	Age  int    `bson:"age"`

	P *Person
}

type Person struct {
	PName string `form:"p"`
}

func TestBind(t *testing.T) {
	var h *Human = &Human{}

	Bind(h, map[string][]string{"name": []string{"adfad"}, "Age": []string{"1234"}, "p": []string{"aaaa"}})

	fmt.Println("=====")
	fmt.Println(h.Name, h.Age, h.PName, h.P.PName)
	fmt.Println(h.CleanData)
}
