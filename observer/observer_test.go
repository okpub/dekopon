package observer

import (
	"fmt"
	"testing"
)

type Data struct {
	Name string
}

type User struct {
	Name string
}

func (u *User) Signup(name string, age int) *Data {
	fmt.Println("Signup name:", name, "age:", age)

	return &Data{Name: name}
}

func (u *User) Signin() {
	fmt.Println("Signin name:", u.Name)
}

func TestObserver(t *testing.T) {
	user := &User{
		Name: "zhangsan",
	}

	Register(user)
	Router("User", "Signup", "abc", 123)

	var b = NewObserver(user)

	b.Call("Signin")

	var obj, _ = GetValue(b.Call("Signup", "123", 432))
	fmt.Println("返回值", obj)
}
