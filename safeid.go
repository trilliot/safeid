package safeid

import (
	"database/sql/driver"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

const (
	stringCodingBase = 62
	stringLen        = 22
	typeSeparator    = "_"
)

// ParseError is returned by FromString and FromUUID when trying to parse a string
// representation that doesn't match ID's prefix or is invalid.
type ParseError string

func (err ParseError) Error() string { return string(err) }

// ID is a unique type-safe identifier backed by a UUID.
// If there is a prefix, it will be separated from the UUID part by an underscore:
// typ_xxxxxxxxxxxxxxxxxxxxxx
//
// A zero-value UUID (00000000-0000-0000-0000-000000000000) will represent a zero-value ID.
type ID[T Prefixer] struct {
	uuid uuid.UUID
}

// IsZero checks if the provided ID is nil or set to its zero-value.
func IsZero[T Prefixer](id *ID[T]) bool {
	if id == nil {
		return true
	}
	return id.uuid == uuid.UUID{}
}

// Must panics if err != nil, otherwise it returns id.
func Must[T Prefixer](id *ID[T], err error) *ID[T] {
	if err != nil {
		panic(err)
	}
	return id
}

// New creates a new type-safe identifier.
// If the Prefixer returns an empty prefix and is not Generic, New panics.
// If an error is returned, a zero-value ID is also returned.
func New[T Prefixer]() (ID[T], error) {
	if prefixOf[T]() == "" && !IsGeneric[T]() {
		panic("empty prefix is not allowed")
	}

	uuidV7, err := uuid.NewV7()
	if err != nil {
		return ID[T]{}, err
	}

	return ID[T]{uuidV7}, nil
}

// FromString parses the string type-safe representation of an ID.
// If the Prefixer returns an empty prefix and is not Generic, FromString panics.
// If the type-safe string representation doesn't max the Prefixer's prefix, FromString
// returns ErrInvalidPrefix.
// If an error is returned, a zero-value ID is also returned.
func FromString[T Prefixer](s string) (ID[T], error) {
	prefix := prefixOf[T]()
	if strings.TrimSpace(prefix) == "" && !IsGeneric[T]() {
		panic("empty prefix is not allowed")
	}

	if s == "" {
		return ID[T]{}, nil
	}

	if prefix != "" {
		var ok bool
		s, ok = strings.CutPrefix(s, prefix+typeSeparator)
		if !ok {
			return ID[T]{}, ParseError("invalid prefix")
		}
	}

	i, ok := new(big.Int).SetString(s, stringCodingBase)
	if !ok {
		return ID[T]{}, ParseError("invalid format")
	}

	// Add leading zeroes if necessary
	ib := i.Bytes()
	b := make([]byte, 16)
	copy(b[16-len(ib):], ib[:])

	return ID[T]{uuid.UUID(b)}, nil
}

// FromUUID parses the string UUID representation of an ID.
// If the Prefixer returns an empty prefix and is not Generic, FromUUID panics.
// If an error is returned, a zero-value ID is also returned.
func FromUUID[T Prefixer](s string) (ID[T], error) {
	if prefixOf[T]() == "" && !IsGeneric[T]() {
		panic("empty prefix is not allowed")
	}

	if s == "" {
		return ID[T]{}, nil
	}

	uuidV7, err := uuid.Parse(s)
	if err != nil {
		return ID[T]{}, ParseError("invalid format")
	}

	return ID[T]{uuidV7}, nil
}

// String returns the type-safe representation of the ID as a string.
func (id ID[T]) String() string {
	prefix := prefixOf[T]()

	i := new(big.Int).SetBytes(id.uuid[:])
	suffix := i.Text(stringCodingBase)
	if len(suffix) < stringLen {
		suffix = strings.Repeat("0", stringLen-len(suffix)) + suffix
	}

	if strings.TrimSpace(prefix) == "" {
		return suffix
	}
	return prefix + typeSeparator + suffix
}

// UUID returns the UUID representation of the ID as a string.
func (id ID[T]) UUID() string {
	return id.uuid.String()
}

// MarshalText satisfies the encoding.TextMarshaler interface by encoding an ID
// to its type-safe string representation.
func (id ID[T]) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

// UnmarshalText satisfies the encoding.TextUnmarshaler interface by decoding a
// type-safe string representation.
func (id *ID[T]) UnmarshalText(data []byte) error {
	v, err := FromString[T](string(data))
	if err == nil {
		*id = v
	}
	return err
}

// Value satisfies the driver.Valuer interface by returning a string representation
// of the UUID.
func (id ID[T]) Value() (driver.Value, error) {
	return id.uuid.Value()
}

// Scan satisfies the sql.Scanner interface by decoding a string representation
// of the UUID.
func (id *ID[T]) Scan(src any) error {
	if prefixOf[T]() == "" && !IsGeneric[T]() {
		panic("empty prefix is not allowed")
	}

	var v uuid.UUID
	err := v.Scan(src)
	if err == nil {
		*id = ID[T]{v}
	}
	return err
}
