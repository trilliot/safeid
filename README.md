# SafeID

#### A human-readable, K-sortable, type-safe at compile-time, UUIDv7-backed globally unique identifiers for Go.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

When dealing with commonly accepted unique identifiers such as UUIDs, 
it is hard to know what they represent at first glance.

SafeIDs are designed to be:

- **human readable:** SafeIDs are prefixed by a _type_, followed by a 22 characters string (e.g.: `user_02Yjy1AYf1ckS6jBZ5zw3G`)
- **K-sortable:** based on the UUIDv7 specification, both UUID or String representation of a SafeID are K-sortable 
- **type-safe:** it is not possible to parse one SafeID type into another when dealing with their string representation
- **compile-time safe:** thanks to generics, your code will not compile if you try to pass one SafeID type into another
- **database efficient:** when stored into a database, SafeID can leverage `uuid` column types

## Installation

```
go get github.com/trilliot/safeid
```

## Usage

You can create either generic (`safeid.Generic`) or custom typed SafeID:
```go
// Create a generic type without prefix (not recommended)
id, err := safeid.New[safeid.Generic]()

// Or create a custom type that satisfies the safeid.Prefixer interface
type User struct {}
func (User) Prefix() string { return "user" }

id, err := safeid.New[User]()
```

An SafeID can be retrieved in its String or UUID form:
```go
id.String() // user_02Yjy1AYf1ckS6jBZ5zw3G
id.UUID()   // 0193508b-e85e-7812-ba4b-91d85495d7bc
```

When dealing with JSON, the type-safe form is (un)marshaled:
```go
obj := map[string]any{
	"id": safeid.Must(safeid.New[User]()),
}
json.Marshal(obj, ...) // {"id": "user_02Yjy1AYf1ckS6jBZ5zw3G"}
```

When passed to or scanned from PostgreSQL, the UUID form is used:
```go
var id safeid.ID[User]
rows.Scan(&id) // 0193508b-e85e-7812-ba4b-91d85495d7bc
```
