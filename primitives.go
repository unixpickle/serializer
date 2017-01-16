package serializer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
)

// Bytes is a serializable version of a byte slice.
type Bytes []byte

func (b Bytes) Serialize() ([]byte, error) {
	return b, nil
}

func (b Bytes) SerializerType() string {
	return "[]byte"
}

// Int is a serializer version of an int.
type Int int

func (i Int) Serialize() ([]byte, error) {
	return []byte(strconv.Itoa(int(i))), nil
}

func (i Int) SerializerType() string {
	return "int"
}

// Float64 is a Serializer for a float64.
type Float64 float64

func (f Float64) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, float64(f))
	return buf.Bytes(), nil
}

func (f Float64) SerializerType() string {
	return "float64"
}

// A Float64Slice is a Serializer for a []float64.
type Float64Slice []float64

func (f Float64Slice) Serialize() ([]byte, error) {
	var w bytes.Buffer
	binary.Write(&w, binary.LittleEndian, uint64(len(f)))
	for _, x := range f {
		binary.Write(&w, binary.LittleEndian, x)
	}
	return w.Bytes(), nil
}

func (f Float64Slice) SerializerType() string {
	return "[]float64"
}

func init() {
	RegisterDeserializer(Bytes(nil).SerializerType(), func(d []byte) (Serializer, error) {
		return Bytes(d), nil
	})
	RegisterDeserializer(Int(0).SerializerType(), func(d []byte) (Serializer, error) {
		num, err := strconv.Atoi(string(d))
		if err != nil {
			return nil, err
		}
		return Int(num), nil
	})
	RegisterDeserializer(Float64(0).SerializerType(), func(d []byte) (Serializer, error) {
		buf := bytes.NewBuffer(d)
		var value float64
		if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
			return nil, errors.New("failed to decode float64: " + err.Error())
		}
		return Float64(value), nil
	})
	RegisterDeserializer(Float64Slice(nil).SerializerType(), func(d []byte) (Serializer, error) {
		reader := bytes.NewBuffer(d)
		var size uint64
		if err := binary.Read(reader, binary.LittleEndian, &size); err != nil {
			return nil, err
		}
		vec := make(Float64Slice, int(size))
		for i := range vec {
			if err := binary.Read(reader, binary.LittleEndian, &vec[i]); err != nil {
				return nil, err
			}
		}
		return vec, nil
	})
}
