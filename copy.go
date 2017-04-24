package serializer

import (
	"fmt"

	"github.com/unixpickle/essentials"
)

// A Copier can produce a deep copy of itself.
type Copier interface {
	Copy() (interface{}, error)
}

// Copy produces a copy of the object.
//
// This uses the Copier interface if possible and falls
// back on using Serializer if necessary.
func Copy(obj interface{}) (copied interface{}, err error) {
	defer essentials.AddCtxTo("copy", &err)
	if copier, ok := obj.(Copier); ok {
		return copier.Copy()
	} else if ser, ok := obj.(Serializer); ok {
		data, err := SerializeWithType(ser)
		if err != nil {
			return nil, err
		}
		return DeserializeWithType(data)
	} else {
		return nil, fmt.Errorf("cannot copy objects of type %T", obj)
	}
}
