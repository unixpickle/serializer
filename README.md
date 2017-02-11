# serializer

Package **serializer** makes it easy to serialize and deserialize objects across instances of a Go application.

To see how to use this package, you can check out the [GoDoc](http://godoc.org/github.com/unixpickle/serializer).

# Why you need it

Go already gives you [encoding/json](https://golang.org/pkg/encoding/json/) and [encoding/gob](https://golang.org/pkg/encoding/gob/), so why do you need this new serialization package? You may not, depending on your use case.

JSON and Gob are great for many situations. JSON is great when you need a cross-platform data format, and Gob is great for sending Go objects over a network. Unfortunately, Gob is not intended to be used as a data storage format; when you use Gob, the data you store depends too much on the layout of the Go structures. JSON, while intended to store data, doesn't work very well with type information. Take this example:

```go
type ActivationFunc interface {
    Eval(x float64) float64
}

type SerializeMe struct {
    X  int
    A  ActivationFunc
    As []ActivationFunc
}
```

The `json` package would succeed at serializing an instance of `SerializeMe`, provided that the `ActivationFunc`s didn't have cyclic structure. However, it would run into trouble trying to deserialize the same instance, since it would have no way to create new instances of `ActivationFunc`.

# How it works

Serializer tries to retain type information while remaining suitable for data storage. To deal with type information, serializer uses type IDs. Any object that implements the `Serializer` interface has a type ID, and this type ID must be registered in a table of decoder functions. This way, any object that can be encoded can also be decoded, since the decoder function knows how to create a new instance of the encoded type.

To keep serialization cross-platform and suitable for data storage, serialization is implemented on a per-type basis. A `Serializer` implements a `Serialize` method which manually encodes the object as binary. The package provides various serialization helpers, but ultimately it is up to an object to decide how it's laid out as data. This makes serializer more cross-platform than Gob in many cases.
