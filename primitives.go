package serializer

import "strconv"

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
}
