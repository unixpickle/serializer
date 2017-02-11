package serializer

import (
	"errors"
	"fmt"
)

func init() {
	var m MyObject
	RegisterTypedDeserializer(m.SerializerType(), DeserializeMyObject)
}

type MyObject struct {
	Num int
	Str string
}

func DeserializeMyObject(d []byte) (*MyObject, error) {
	var num Int
	var str Bytes
	if err := DeserializeAny(d, &num, &str); err != nil {
		return nil, errors.New("deserialize MyObject: " + err.Error())
	}
	return &MyObject{Num: int(num), Str: string(str)}, nil
}

func (m *MyObject) SerializerType() string {
	return "github.com/unixpickle/serializer.MyObject"
}

func (m *MyObject) Serialize() ([]byte, error) {
	return SerializeAny(Int(m.Num), Bytes([]byte(m.Str)))
}

func Example() {
	obj := &MyObject{15, "hello, world"}
	data, _ := SerializeAny(obj)

	var newObj *MyObject
	if err := DeserializeAny(data, &newObj); err != nil {
		panic(err)
	}

	fmt.Println(newObj.Num, newObj.Str)

	// Output: 15 hello, world
}
