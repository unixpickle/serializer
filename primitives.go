package serializer

// Bytes is a serializable version of a byte slice.
type Bytes []byte

func (b Bytes) Serialize() ([]byte, error) {
	return b, nil
}

func (b Bytes) SerializerType() string {
	return "[]byte"
}

func init() {
	RegisterDeserializer(Bytes(nil).SerializerType(), func(d []byte) (Serializer, error) {
		return Bytes(d), nil
	})
}
