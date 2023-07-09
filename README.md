## Project description
Tool for working with structure fields.

## Install
```bash
go get github.com/ruauka/attrs-go@v1.0.9
```

## Usage
<ins>**GetAttr**</ins> - get struct field value. 

Args:
 - obj: struct - value param. Fields can be ptr or value.
 - fieldName: value param.
```go
func GetAttr(obj interface{}, fieldName string) (interface{}, error)
```

```go
type User struct {
    Username string
}

user := User{Username: "username value"}

value, err := attrs_go.GetAttr(user, "Username")
fmt.Println(value) // username value
```

<ins>**SetAttr**</ins>  - set new value on structure field.

Args:
- obj: struct - ptr param. Struct fields can be ptr or value.
- fieldName, newValue: value param.

```go
func SetAttr(obj interface{}, fieldName string, newValue interface{}) error
```

```go
type User struct {
    Username string
}

u := &User{Username: "username value"}

if err := attrs_go.SetAttr(u, "Username", "new username value"); err != nil {
    fmt.Println(err)
}

fmt.Println(u.Username) // new username value
```

<ins>**SetStructAttrs**</ins>  - updates current structure fields with the values of the new structure fields.

Args:
- curObj: struct - ptr param. Struct fields can be ptr or value.
- newObj: struct - value param. Struct fields can be ptr or value.

```go
func SetStructAttrs(curObj, newObj interface{}) error
```

```go
type User struct {
    Username string // will change by pte
    Age      int    // will change by value
    Married  bool   // will be the same
}

type NewUser struct {
    Username *string `json:"username"`
    Age      int     `json:"age"`
    Married  *bool   `json:"married"` // nil
}

user := &User{
    Username: "username",
    Age:      30,
    Married: true
}

newUserName := "new_username"
newUser := NewUser{
    Username: &newUserName,
    Age:      35,
    Married:  nil,
}

fmt.Printf("%s, %d, %v\n", user.Username, user.Age, user.Married) // username, 30, true

if err := attrs_go.SetStructAttrs(user, newUser); err != nil {
    fmt.Println(err)
}

fmt.Printf("%s, %d, %v\n", user.Username, user.Age, user.Married) // new_username, 35, true
```

<ins>**RoundUpFloatStruct**</ins>  - round up float struct fields to certain precision.

Args:
- obj: struct - ptr param. Struct fields can be value, not ptr.
- precision: round to.

Constraint:

- simple floats
- array and slice

```go
func RoundUpFloatStruct(obj interface{}, precision int) error
```

```go
type Foo struct {
    Field1 float32
    Field2 float64
    Field3 []float32
    Field4 []float64
    Field5 [3]float32
    Field6 [3]float64
    Field7 int // will be the same
    Field8 string // will be the same
}

foo := &Foo{
    Field1: 1.1111,
    Field2: 2.2222,
    Field3: []float32{1.1111, 2.2222, 3.3333},
    Field4: []float64{4.4444, 5.5555, 7.7777},
    Field5: [3]float32{1.1111, 2.2222, 3.3333},
    Field6: [3]float64{4.4444, 5.5555, 7.7777},
    Field7: 7,
    Field8: "field8",
}

fmt.Printf("%+v\n", *foo)
// {
//Field1:1.1111 Field2:2.2222
//Field3:[1.1111 2.2222 3.3333] Field4:[4.4444 5.5555 7.7777]
//Field5:[1.1111 2.2222 3.3333] Field6:[4.4444 5.5555 7.7777]
//Field7:7 Field8:field8
//}

if err := attrs_go.RoundUpFloatStruct(foo, 3); err != nil {
    fmt.Println(err)
}

fmt.Printf("%+v", *foo)
// {
//Field1:1.112 Field2:2.223
//Field3:[1.112 2.223 3.334] Field4:[4.445 5.556 7.778]
//Field5:[1.112 2.223 3.334] Field6:[4.445 5.556 7.778]
//Field7:7 Field8:field8
//}
```