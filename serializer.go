package serializer

// A Serializer is any object which can be serialized.
type Serializer interface {
	// Serialize serializes this object.
	Serialize() ([]byte, error)

	// SerializerType returns the unique type ID
	// for this type, which should be registered
	// in the decoder table.
	SerializerType() string
}
