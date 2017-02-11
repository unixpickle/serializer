// Package serializer provides a type-aware data storage
// system.
// It can be used to transfer objects over a network, save
// objects to disk, etc.
//
// Serializer tries to retain type information while
// remaining suitable for data storage.
// To deal with type information, it uses type IDs.
// Any object that implements the Serializer interface has
// a type ID, and this type ID must be registered in a
// table of decoder functions.
// This way, any object that can be encoded can also be
// decoded, since the decoder function knows how to create
// a new instance of the encoded type.
//
// To keep serialization cross-platform and suitable for
// data storage, serialization is implemented on a
// per-type basis.
// A Serializer implements a Serialize method which
// manually encodes the object as binary.
// The package provides various serialization helpers, but
// ultimately it is up to an object to decide how it's
// laid out as data.
// This makes serializer more cross-platform than Gob in
// many cases.
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
