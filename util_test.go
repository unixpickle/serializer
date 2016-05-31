package serializer

import (
	"encoding/json"
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
