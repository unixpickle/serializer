package serializer

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestPrimitives(t *testing.T) {
	objects := []interface{}{
		Bytes([]byte("hello, world")),
		Int(1337),
		Int(-(1 << 30)),
		Float64(3.1415),
		Float32(3.1415),
		Float64Slice([]float64{1, 0.5, -5.15e20}),
		Float32Slice([]float32{-1337, 0.5, 5e10, -5e10}),
		Float64Slice([]float64{}),
		Float32Slice([]float32{}),
		String("hello, world"),
		Bool(true),
		Bool(false),
		IntSlice([]int{1, 2, 3}),
		Int32(13371337),
		Int64(133713371337),
		Int32Slice([]int32{1, 2, 15}),
		Int64Slice([]int64{1, 2, 133713371337}),
	}
	data, err := SerializeAny(objects...)
	if err != nil {
		t.Fatal(err)
	}
	var obj1 Bytes
	var obj2, obj3 Int
	var obj4 Float64
	var obj5 Float32
	var obj6, obj8 Float64Slice
	var obj7, obj9 Float32Slice
	var obj10 String
	var obj11, obj12 Bool
	var obj13 IntSlice
	var obj14 Int32
	var obj15 Int64
	var obj16 Int32Slice
	var obj17 Int64Slice
	err = DeserializeAny(data, &obj1, &obj2, &obj3, &obj4, &obj5, &obj6, &obj7, &obj8,
		&obj9, &obj10, &obj11, &obj12, &obj13, &obj14, &obj15, &obj16, &obj17)
	if err != nil {
		t.Fatal(err)
	}
	newObjs := []interface{}{obj1, obj2, obj3, obj4, obj5, obj6, obj7, obj8,
		obj9, obj10, obj11, obj12, obj13, obj14, obj15, obj16, obj17}
	for i, x := range objects {
		if !reflect.DeepEqual(x, newObjs[i]) {
			t.Errorf("object %d: expected %v but got %v", i, x, newObjs[i])
		}
	}
}

func BenchmarkFloat32Serialize(b *testing.B) {
	buf := make([]float32, 1000000)
	for i := range buf {
		buf[i] = rand.Float32()
	}
	for i := 0; i < b.N; i++ {
		Float32Slice(buf).Serialize()
	}
}

func BenchmarkFloat32Deserilaize(b *testing.B) {
	buf := make([]float32, 1000000)
	for i := range buf {
		buf[i] = rand.Float32()
	}
	data, _ := Float32Slice(buf).Serialize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeserializeFloat32Slice(data)
	}
}

func BenchmarkFloat64Serialize(b *testing.B) {
	buf := make([]float64, 1000000)
	for i := range buf {
		buf[i] = rand.Float64()
	}
	for i := 0; i < b.N; i++ {
		Float64Slice(buf).Serialize()
	}
}

func BenchmarkFloat64Deserilaize(b *testing.B) {
	buf := make([]float64, 1000000)
	for i := range buf {
		buf[i] = rand.Float64()
	}
	data, _ := Float64Slice(buf).Serialize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeserializeFloat64Slice(data)
	}
}

func BenchmarkIntSerialize(b *testing.B) {
	buf := make([]int, 1000000)
	for i := range buf {
		buf[i] = rand.Int()
	}
	for i := 0; i < b.N; i++ {
		IntSlice(buf).Serialize()
	}
}

func BenchmarkIntDeserialize(b *testing.B) {
	buf := make([]int, 1000000)
	for i := range buf {
		buf[i] = rand.Int()
	}
	data, _ := IntSlice(buf).Serialize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeserializeIntSlice(data)
	}
}
