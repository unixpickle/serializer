# serializer

Package **serializer** makes it easy to serialize and deserialize objects across instances of a Go application.

To see how to use this package, you can check out the [GoDoc](http://godoc.org/github.com/unixpickle/serializer).

# Why you need it

Go already gives you [encoding/json](https://golang.org/pkg/encoding/json/) and [encoding/gob](https://golang.org/pkg/encoding/gob/), so why do you need this new serialization package? You may not, depending on your use case.

The `json` and `gob` packages are great for many situations. Unfortunately, they can't help when your objects have references to interfaces. For instance, suppose you have this situation:

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

The standard Go serialization packages (`json` and `gob`) would probably succeed at serializing an instance of `SerializeMe`, provided that the `ActivationFunc`s didn't have cyclic structure. However, they would run into trouble trying to deserialize the same instance, since they'd have no way to create new instances of `ActivationFunc`.

# How it works

All *serializer* adds to `json` or `gob` is type information. Any object that implements the `Serializer` interface has a type ID, and this type ID must be registered in a table of decoder functions. This way, any object that can be encoded can also be decoded, since the decoder function knows how to create a new instance of the encoded type.
