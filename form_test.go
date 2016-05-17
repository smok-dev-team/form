package form

import (
	"fmt"
	"testing"
)

type Human struct {
	Person
	CleanData map[string]interface{}
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

type Person struct {
	PName string `form:"p"`
}

func TestBind(t *testing.T) {
	var h Human
	h.CleanData = make(map[string]interface{})

	Bind(&h, map[string][]string{"Name": []string{"adfad"}, "Age": []string{"1234"}, "p": []string{"aaaa"}})

	fmt.Println("=====")
	fmt.Println(h.Name, h.Age, h.PName)
	fmt.Println(h.CleanData)
}
