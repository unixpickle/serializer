package serializer

import (
	"reflect"
	"sync"
)

var deserializersLock sync.RWMutex
var deserializers = map[string]Deserializer{}

// A Deserializer is a function which can deserialize
// a certain type of object.
type Deserializer func(d []byte) (Serializer, error)

// GetDeserializer returns the Deserializer that is
// currently registered for the given type ID.
// This returns nil if no Deserializer is registered.
//
// All routines which manage the deserializer table
// are safe to call concurrently.
func GetDeserializer(typeID string) Deserializer {
	deserializersLock.RLock()
	defer deserializersLock.RUnlock()
	return deserializers[typeID]
}

// UpdateDeserializer adds or changes the Deserializer
// for a given type ID.
//
// Passing a nil Deserializer will completely remove
// the type ID from the table, allowing
// RegisterDeserializer to be called again.
//
// All routines which manage the deserializer table
// are safe to call concurrently.
func UpdateDeserializer(typeID string, d Deserializer) {
	deserializersLock.Lock()
	defer deserializersLock.Unlock()
	if d == nil {
		delete(deserializers, typeID)
	} else {
		deserializers[typeID] = d
	}
}

// RegisterDeserializer is like UpdateDeserializer,
// but it panics if the type ID is already in use.
//
// All routines which manage the deserializer table
// are safe to call concurrently.
func RegisterDeserializer(typeID string, d Deserializer) {
	deserializersLock.Lock()
	defer deserializersLock.Unlock()
	if _, ok := deserializers[typeID]; ok {
		panic("type ID already in use: " + typeID)
	}
	deserializers[typeID] = d
}

// RegisterTypedDeserializer is like RegisterDeserializer,
// but instead of taking a Deserializer, it converts a
// function into a Deserializer by casting its first return
// value to a Serializer.
//
// For instance, you might have a method like this:
//
//     func DeserializeMyType(d []byte) (*MyType, error) {
//         ...
//     }
//
// Technically, DeserializeMyType is not a Deserializer
// since it doesn't return a Serializer, it returns *MyType.
// Still, it is clear how it could be turned into a
// Deserializer, provided that *MyType implements Serializer.
// The job of RegisterTypedDeserializer is to use a method
// like DeserializeMyType as a deserializer.
func RegisterTypedDeserializer(typeID string, f interface{}) {
	val := reflect.ValueOf(f)
	RegisterDeserializer(typeID, func(d []byte) (Serializer, error) {
		res := val.Call([]reflect.Value{reflect.ValueOf(d)})
		if res[1].IsNil() {
			return res[0].Interface().(Serializer), nil
		} else {
			return nil, res[1].Interface().(error)
		}
	})
}
