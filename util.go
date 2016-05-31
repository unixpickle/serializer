package serializer

import (
	"bytes"
	"encoding/binary"
	"errors"
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
