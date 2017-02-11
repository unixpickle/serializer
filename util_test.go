package serializer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

type demoType1 struct {
	X int
	Y int
}

func (d *demoType1) Serialize() ([]byte, error) {
	return json.Marshal(d)
}

func (d *demoType1) SerializerType() string {
	return "demoType1"
}

type demoType2 struct {
	Z float64
	K float64
}

func (d *demoType2) Serialize() ([]byte, error) {
	return json.Marshal(d)
}

func (d *demoType2) SerializerType() string {
	return "demoType2"
}

func deserializeDemoType1(d []byte) (Serializer, error) {
	var x demoType1
	if err := json.Unmarshal(d, &x); err != nil {
		return nil, err
	}
	return &x, nil
}

func deserializeDemoType2(d []byte) (Serializer, error) {
	var x demoType2
	if err := json.Unmarshal(d, &x); err != nil {
		return nil, err
	}
	return &x, nil
}

func TestSerializeSlice(t *testing.T) {
	UpdateDeserializer("demoType1", deserializeDemoType1)
	UpdateDeserializer("demoType2", deserializeDemoType2)

	slice := []Serializer{
		&demoType1{X: 1, Y: -100},
		&demoType2{Z: 3.14, K: 3.41},
		&demoType2{K: 5.3},
		&demoType1{Y: 50},
	}
	serialized, err := SerializeSlice(slice)
	if err != nil {
		t.Fatal(err)
	}
	deserialized, err := DeserializeSlice(serialized)
	if err != nil {
		t.Fatal(err)
	}
	if len(deserialized) != len(slice) {
		t.Fatal("array lengths do not match: got", len(deserialized), "expected", len(slice))
	}
	if *(slice[0].(*demoType1)) != *(deserialized[0].(*demoType1)) {
		t.Error("element 0 does not match")
	}
	if *(slice[1].(*demoType2)) != *(deserialized[1].(*demoType2)) {
		t.Error("element 1 does not match")
	}
	if *(slice[2].(*demoType2)) != *(deserialized[2].(*demoType2)) {
		t.Error("element 2 does not match")
	}
	if *(slice[3].(*demoType1)) != *(deserialized[3].(*demoType1)) {
		t.Error("element 3 does not match")
	}
}

func TestSerializeAny(t *testing.T) {
	var obj Int
	var obj1 interface{}

	obj = 7
	obj1 = Bytes([]byte("hello"))

	if data, err := SerializeAny(obj, obj1); err != nil {
		t.Error(err)
	} else {
		obj = 3
		obj1 = "hello"
		if err := DeserializeAny(data, &obj, &obj1); err != nil {
			t.Error(err)
		}
		if obj != 7 {
			t.Errorf("expected %d got %d", 7, obj)
		}
		if b, ok := obj1.(Bytes); !ok {
			t.Errorf("expected Bytes got %T", obj1)
		} else if !bytes.Equal(b, []byte("hello")) {
			t.Errorf("expected %v got %v", []byte("hello"), b)
		}

		var obj3 string
		if err := DeserializeAny(data, &obj, &obj3); err == nil {
			t.Errorf("expecting error")
		}
		if err := DeserializeAny(data, &obj1, &obj); err == nil {
			t.Errorf("expceting error")
		}
	}
}

func ExampleSerializeAny() {
	obj1 := Int(15)
	obj2 := Float64(3.14)
	data, _ := SerializeAny(obj1, obj2)

	var out1 Int
	var out2 Float64
	DeserializeAny(data, &out1, &out2)

	fmt.Println(out1, out2)

	// Output: 15 3.14
}
