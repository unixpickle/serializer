package serializer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"

	"github.com/unixpickle/essentials"
)

func init() {
	RegisterTypedDeserializer(Bytes(nil).SerializerType(), DeserializeBytes)
	RegisterTypedDeserializer(String("").SerializerType(), DeserializeString)
	RegisterTypedDeserializer(Int(0).SerializerType(), DeserializeInt)
	RegisterTypedDeserializer(IntSlice(nil).SerializerType(), DeserializeIntSlice)
	RegisterTypedDeserializer(Float64(0).SerializerType(), DeserializeFloat64)
	RegisterTypedDeserializer(Float32(0).SerializerType(), DeserializeFloat32)
	RegisterTypedDeserializer(Float64Slice(nil).SerializerType(), DeserializeFloat64Slice)
	RegisterTypedDeserializer(Float32Slice(nil).SerializerType(), DeserializeFloat32Slice)
	RegisterTypedDeserializer(Bool(false).SerializerType(), DeserializeBool)
}

// Bytes is a Serializer wrapper for []byte.
type Bytes []byte

// DeserializeBytes deserializes a Bytes.
func DeserializeBytes(d []byte) (Bytes, error) {
	return d, nil
}

// Serialize serializes the object.
func (b Bytes) Serialize() ([]byte, error) {
	return b, nil
}

// SerializerType returns the unique ID used to serialize
// Bytes.
func (b Bytes) SerializerType() string {
	return "[]byte"
}

// String is a Serializer wrapper for string.
type String string

// DeserializeString deserializes a String.
func DeserializeString(d []byte) (String, error) {
	return String(d), nil
}

// Serialize serializes the object.
func (s String) Serialize() ([]byte, error) {
	return []byte(s), nil
}

// SerializerType returns the unique ID used to serialize
// a String.
func (s String) SerializerType() string {
	return "string"
}

// Int is a Serializer wrapper for an int.
type Int int

// DeserializeInt deserialize an Int.
func DeserializeInt(d []byte) (Int, error) {
	num, err := strconv.Atoi(string(d))
	if err != nil {
		return 0, essentials.AddCtx("deserialize int", err)
	}
	return Int(num), nil
}

// Serialize serializes the object.
func (i Int) Serialize() ([]byte, error) {
	return []byte(strconv.Itoa(int(i))), nil
}

// SerializerType returns the unique ID used to serialize
// an Int.
func (i Int) SerializerType() string {
	return "int"
}

// IntSlice is a Serializer wrapper for an []int.
type IntSlice []int

// DeserializeIntSlice deserializes an IntSlice.
func DeserializeIntSlice(d []byte) (slice IntSlice, err error) {
	defer essentials.AddCtxTo("deserialize IntSlice", &err)

	reader := bytes.NewBuffer(d)

	var size uint64
	if err := binary.Read(reader, binary.LittleEndian, &size); err != nil {
		return nil, err
	}
	vec := make([]int64, int(size))
	if err := binary.Read(reader, binary.LittleEndian, vec); err != nil {
		return nil, err
	}

	res := make([]int, len(vec))
	for i, x := range vec {
		res[i] = int(x)
	}

	return res, nil
}

// Serialize serializes the object.
func (i IntSlice) Serialize() ([]byte, error) {
	ints64 := make([]int64, len(i))
	for j, x := range i {
		ints64[j] = int64(x)
	}
	var w bytes.Buffer
	binary.Write(&w, binary.LittleEndian, uint64(len(i)))
	binary.Write(&w, binary.LittleEndian, ints64)
	return w.Bytes(), nil
}

// SerializerType returns the unique ID used to serialize
// an IntSlice.
func (i IntSlice) SerializerType() string {
	return "[]int"
}

// Float64 is a Serializer for a float64.
type Float64 float64

// DeserializeFloat64 deserializes a Float64.
func DeserializeFloat64(d []byte) (Float64, error) {
	buf := bytes.NewBuffer(d)
	var value float64
	if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
		return 0, essentials.AddCtx("deserialize float64", err)
	}
	return Float64(value), nil
}

// Serialize serializes the object.
func (f Float64) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, float64(f))
	return buf.Bytes(), nil
}

// SerializerType returns the unique ID used to serialize
// a Float64.
func (f Float64) SerializerType() string {
	return "float64"
}

// Float32 is a Serializer for a float32.
type Float32 float32

// DeserializeFloat32 deserializes a Float32.
func DeserializeFloat32(d []byte) (Float32, error) {
	buf := bytes.NewBuffer(d)
	var value float32
	if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
		return 0, essentials.AddCtx("deserialize float32", err)
	}
	return Float32(value), nil
}

// Serialize serializes the object.
func (f Float32) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, float32(f))
	return buf.Bytes(), nil
}

// SerializerType returns the unique ID used to serialize
// a Float64.
func (f Float32) SerializerType() string {
	return "float32"
}

// A Float64Slice is a Serializer for a []float64.
type Float64Slice []float64

// DeserializeFloat64Slice deserializes a Float64Slice.
func DeserializeFloat64Slice(d []byte) (Float64Slice, error) {
	reader := bytes.NewBuffer(d)
	var size uint64
	if err := binary.Read(reader, binary.LittleEndian, &size); err != nil {
		return nil, essentials.AddCtx("deserialize []float64", err)
	}
	vec := make([]float64, int(size))
	if err := binary.Read(reader, binary.LittleEndian, vec); err != nil {
		return nil, essentials.AddCtx("deserialize []float64", err)
	}
	return vec, nil
}

// Serialize serializes the object.
func (f Float64Slice) Serialize() ([]byte, error) {
	var w bytes.Buffer
	binary.Write(&w, binary.LittleEndian, uint64(len(f)))
	binary.Write(&w, binary.LittleEndian, []float64(f))
	return w.Bytes(), nil
}

// SerializerType returns the unique ID used to serialize
// a Float64Slice.
func (f Float64Slice) SerializerType() string {
	return "[]float64"
}

// A Float32Slice is a Serializer for a []float32.
type Float32Slice []float32

// DeserializeFloat32Slice deserializes a Float32Slice.
func DeserializeFloat32Slice(d []byte) (Float32Slice, error) {
	reader := bytes.NewBuffer(d)
	var size uint64
	if err := binary.Read(reader, binary.LittleEndian, &size); err != nil {
		return nil, essentials.AddCtx("deserialize []float32", err)
	}
	vec := make([]float32, int(size))
	if err := binary.Read(reader, binary.LittleEndian, vec); err != nil {
		return nil, essentials.AddCtx("deserialize []float32", err)
	}
	return vec, nil
}

// Serialize serializes the object.
func (f Float32Slice) Serialize() ([]byte, error) {
	var w bytes.Buffer
	binary.Write(&w, binary.LittleEndian, uint64(len(f)))
	binary.Write(&w, binary.LittleEndian, []float32(f))
	return w.Bytes(), nil
}

// SerializerType returns the unique ID used to serialize
// a Float32Slice.
func (f Float32Slice) SerializerType() string {
	return "[]float32"
}

// A Bool is a Serializer for a bool.
type Bool bool

// DeserializeBool deserializes a Bool.
func DeserializeBool(d []byte) (Bool, error) {
	if len(d) != 1 {
		return false, errors.New("deserialize bool: invalid length")
	}
	if d[0] == 0 {
		return false, nil
	} else if d[0] == 1 {
		return true, nil
	} else {
		return false, errors.New("deserialize bool: invalid value")
	}
}

// Serialize serializes the object.
func (b Bool) Serialize() ([]byte, error) {
	if b {
		return []byte{1}, nil
	} else {
		return []byte{0}, nil
	}
}

// SerializerType returns the unique ID used to serialize
// a Bool.
func (b Bool) SerializerType() string {
	return "bool"
}
