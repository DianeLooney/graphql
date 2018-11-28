package resolver

// Args contains the arguments passed to a Query, Mutation, or Field resolver
type Args map[string]interface{}

// Query represents a query resolver
//   args is the (potentially nil) map of arguments
//   result is the return value, can be nil to express a null value
//   err is the
type Query interface {
	Resolve(args Args) (result interface{}, err error)
}

// Scalar represents a value that can be resolved as a scalar
type Scalar func(sc interface{}) (result interface{}, err error)

// Object represents a map
type Object interface {
	Resolve(field string, args Args) (result interface{}, err error)
}

// Array represents an array
type Array interface {
	Len() int
	Get(i int) (result interface{}, err error)
}
