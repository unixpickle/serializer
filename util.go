package serializer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
)

var helperByteOrder = binary.LittleEndian

var (
	ErrBufferUnderflow = errors.New("buffer underflow")
	ErrResidualData    = errors.New("residual data")
)

// SerializeWithType returns a binary value that
// includes both the type and serialized data for
// the given Serializer.
// This is meant to be used with DeserializeWithType.
func SerializeWithType(s Serializer) ([]byte, error) {
	data, err := s.Serialize()
	if err != nil {
		return nil, err
	}
	typeData := []byte(s.SerializerType())

	res := make([]byte, len(data)+len(typeData)+4)
	copy(res[len(typeData)+4:], data)
	copy(res[4:], typeData)
	helperByteOrder.PutUint32(res, uint32(len(typeData)))

	return res, nil
}

// DeserializeWithType performs the inverse of
// SerializeWithType, first decoding the type ID and
// then using that type ID to decode the object.
func DeserializeWithType(d []byte) (Serializer, error) {
	if len(d) < 4 {
		return nil, ErrBufferUnderflow
	}

	size := int(helperByteOrder.Uint32(d))
	if size+4 > len(d) {
		return nil, ErrBufferUnderflow
	}
	typeID := string(d[4 : size+4])

	deserializer := GetDeserializer(typeID)
	if deserializer == nil {
		return nil, errors.New("unregistered type ID: " + typeID)
	}

	return deserializer(d[4+size:])
}

// SerializeSlice serializes a slice of Serializers,
// storing the size and type ID of each element.
// This is meant to be used in conjunction with
// DeserializeSlice.
func SerializeSlice(s []Serializer) ([]byte, error) {
	var res bytes.Buffer

	for _, x := range s {
		serialized, err := SerializeWithType(x)
		if err != nil {
			return nil, err
		}
		binary.Write(&res, helperByteOrder, uint64(len(serialized)))
		res.Write(serialized)
	}

	return res.Bytes(), nil
}

// DeserializeSlice does the inverse of SerializeSlice.
func DeserializeSlice(d []byte) ([]Serializer, error) {
	buf := bytes.NewBuffer(d)
	var res []Serializer

	for buf.Len() >= 8 {
		var nextLen64 uint64
		binary.Read(buf, helperByteOrder, &nextLen64)
		nextLen := int(nextLen64)
		if nextLen > buf.Len() {
			return nil, ErrBufferUnderflow
		}

		nextData := make([]byte, nextLen)
		buf.Read(nextData)

		obj, err := DeserializeWithType(nextData)
		if err != nil {
			return nil, err
		}
		res = append(res, obj)
	}

	if buf.Len() != 0 {
		return nil, ErrResidualData
	}

	return res, nil
}

// SerializeAny attempts to serialize the objects.
// It fails if any of the objects are not Serializers.
func SerializeAny(obj ...interface{}) ([]byte, error) {
	s := make([]Serializer, len(obj))
	for i, x := range obj {
		var ok bool
		s[i], ok = x.(Serializer)
		if !ok {
			return nil, fmt.Errorf("not a Serializer: %T", x)
		}
	}
	return SerializeSlice(s)
}

// DeserializeAny attempts to reverse the process done by
// SerializeAny.
// It takes pointers to the variables into which the
// objects should be deserialized.
func DeserializeAny(data []byte, out ...interface{}) error {
	slice, err := DeserializeSlice(data)
	if err != nil {
		return err
	}
	if len(slice) != len(out) {
		return fmt.Errorf("have %d destinations but %d decoded objects",
			len(out), len(slice))
	}
	for i, obj := range slice {
		val := reflect.ValueOf(obj)
		destVal := reflect.ValueOf(out[i])
		if destVal.Kind() != reflect.Ptr {
			return fmt.Errorf("element %d: expected pointer but got %T",
				i, out[i])
		}
		if !val.Type().AssignableTo(destVal.Type().Elem()) {
			return fmt.Errorf("element %d: expecting %s but decoded %T",
				i, destVal.Type().Elem(), obj)
		}
		destVal.Elem().Set(val)
	}
	return nil
}

// SaveAny writes the given objects to a file.
// It is like using SerializeAny and writing the results
// to a file afterward.
func SaveAny(path string, obj ...interface{}) error {
	enc, err := SerializeAny(obj...)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, enc, 0755)
}

// LoadAny loads the given objects from a file.
// It is like using DeserializeAny, but first reading the
// data from a file.
func LoadAny(path string, objOut ...interface{}) error {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return DeserializeAny(contents, objOut...)
}
