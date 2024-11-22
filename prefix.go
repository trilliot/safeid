package safeid

// Prefixer is the interface that defines the object type of ID.
// You should implement this interface to define custom type-safe identifiers:
// type ObjectPrefix struct{}
// func (ObjectPrefix) Prefix() string { return "obj" }
type Prefixer interface {
	Prefix() string
}

// Generic is a special identifier type that do not represent a specific object.
// It can be used in cases where there is no need to differentiate identifiers.
// An Generic-type ID cannot be converted to or from another type-safe ID.
type Generic struct{}

// Prefix satisfies the Prefixer interface by returning an empty string.
func (Generic) Prefix() string { return "" }

// IsGeneric checks if the identifier's type is Generic.
func IsGeneric[T Prefixer]() bool {
	var p T
	switch any(p).(type) {
	case Generic:
		return true
	default:
		return false
	}
}

func prefixOf[T Prefixer]() string {
	var p T
	return p.Prefix()
}
