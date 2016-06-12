## FORM
一个用于将 HTTP 请求参数与 Struct 绑定的组件。

通过反射机制，将 HTTP 请求的参数与相关的 Struct 进行绑定，使用很简单：

```
import (
	"fmt"
	"net/http"
	"github.com/smartwalle/form"
)

type Human struct {
	Name string `form:"name"`
	Age  int    `form:"age"`
}

type Student struct {
	Human
	Number int `form:"number"`
}

var formData = map[string][]string{"name": []string{"Yangfeng"}, "age": []string{"12"}, "number": []string{"9"}}

// 绑定
var s *Student
var err = form.Bind(formData, &s)
if err != nil {
	fmt.Println("绑定失败")
}

// 绑定 http.Request, 这样就可以很方便地集成到 Web 项目中。
var s *Student
var err = form.BindWithRequest(request, &s)
if err != nil {
	fmt.Println("绑定失败")
}
```

#### CleanedData
对于一些特殊的需求，比如某个接口要修改用户资料信息，用户资料包含姓名、年龄、出生年月等等。常见的情况就是，不管客户端修改多少个字段的信息，都会将该用户的所有信息都提交一次，这样做有点不科学，浪费资源。但是不这么做，服务器端又不能很方便的知道客户端要修改的字段信息（当然也可以单独用一个参数来传输这次有变化的字段列表），因为对于基本数据类型，如果没有设置值，其都会有一个默认值，比如: string 类型的默认值为空字符串，int 类型的默认值为 0。但是这些默认值不足以用于判断客户端是否需要修改该字段的值。

为了解决这个问题，本 Form 组件提供了一个属性 —— CleanedData，用于存储本次有变化的字段及其值。

CleanedData 的类型为 map[string]interface{}, 其不是一个必须属性，需要由开发者自行在 Struct 中声明，如果没有声明，则没有。如果有声明，在绑定数据的时候，会自动填充相关的信息。其中 key 为 Struct 属性的 form tag，value 为 Struct 对应属性的值。

```
type Human struct {
	CleanedData map[string]interface{}
	Name string `form:"name"`
	Age  int
}

var h Human
Bind(...)

h.CleanedData 包含 name 和 Age 两个key。前提是绑定的数据源也包含这两个字段。

```

#### 默认值
为结构体添加添加 Default+属性名 的方法名，并且返回和该属性相同类型的数据，即可为该属性添加默认值。

例如：

```
func (this Human) DefaultAge() int {
	return 18
}
```

如上所示，如果绑定的数据源中没有找到 Age 相关的字段，则会将 Age 属性的值初始化为 DefaultAge 方法返回的值。